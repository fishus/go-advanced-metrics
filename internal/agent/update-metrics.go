package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/fishus/go-advanced-metrics/internal/secure"
	"github.com/fishus/go-advanced-metrics/internal/storage"
)

func CollectAndPostMetrics(ctx context.Context) {
	dataCh := collectMetricsAtIntervals(ctx)
	postMetricsAtIntervals(ctx, dataCh)
}

// collectMetricsAtIntervals collects metrics every {options.pollInterval} seconds
func collectMetricsAtIntervals(ctx context.Context) chan *storage.MemStorage {
	dataCh := make(chan *storage.MemStorage, 10)

	wgAgent.Add(1)

	go func() {
		defer close(dataCh)
		defer wgAgent.Done()

		ticker := time.NewTicker(Config.PollInterval())
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				data := collectMetrics(ctx)
				if data != nil {
					dataCh <- data
				}
			}
		}
	}()

	return dataCh
}

func collectMetrics(ctx context.Context) *storage.MemStorage {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var mRuntime *storage.MemStorage
	go func() {
		mRuntime = collector.CollectRuntimeMetrics(ctx)
		wg.Done()
	}()

	var mPs *storage.MemStorage
	go func() {
		mPs = collector.CollectPsMetrics(ctx)
		wg.Done()
	}()

	wg.Wait()

	ms := make([]*storage.MemStorage, 0, 2)

	if mRuntime != nil {
		ms = append(ms, mRuntime)
	}

	if mPs != nil {
		ms = append(ms, mPs)
	}

	return combineMetrics(ctx, ms...)
}

func combineMetrics(ctx context.Context, ms ...*storage.MemStorage) *storage.MemStorage {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	data := storage.NewMemStorage()

	for _, m := range ms {
		for _, g := range m.Gauges() {
			_ = data.SetGauge(g.Name(), g.Value())
		}
		for _, c := range m.Counters() {
			_ = data.AddCounter(c.Name(), c.Value())
		}
	}

	if len(data.Gauges()) == 0 && len(data.Counters()) == 0 {
		return nil
	}

	return data
}

// postMetricsAtIntervals posts collected metrics every {options.reportInterval} seconds
func postMetricsAtIntervals(ctx context.Context, dataCh <-chan *storage.MemStorage) {
	dataBuf := make([]*storage.MemStorage, 0)
	workerCh := make(chan *storage.MemStorage, Config.RateLimit())
	defer close(workerCh)

	wgAgent.Add(int(Config.RateLimit()))

	for w := 1; w <= int(Config.RateLimit()); w++ {
		go workerPostMetrics(ctx, workerCh)
	}

	ticker := time.NewTicker(Config.ReportInterval())
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-dataCh:
			if data != nil {
				dataBuf = append(dataBuf, data)
			}
		case <-ticker.C:
			for _, data := range dataBuf {
				workerCh <- data
			}
			dataBuf = dataBuf[:0]
		}
	}
}

// workerPostMetrics posts collected metrics
func workerPostMetrics(ctx context.Context, dataCh <-chan *storage.MemStorage) {
	defer wgAgent.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	client := resty.New().SetBaseURL("http://" + Config.ServerAddr())
	logger.Log.Info("Running agent worker", logger.String("address", Config.ServerAddr()), logger.String("event", "start agent worker"))

	gz, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		logger.Log.Panic(err.Error())
		return
	}

	for data := range dataCh {
		err := retryPostMetrics(ctx, client, gz, data)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

func packMetricsIntoBatch(data *storage.MemStorage) []metrics.Metrics {
	batch := make([]metrics.Metrics, 0)

	for _, counter := range data.Counters() {
		batch = append(batch, metrics.NewCounterMetric(counter.Name()).SetDelta(counter.Value()))
	}

	for _, gauge := range data.Gauges() {
		batch = append(batch, metrics.NewGaugeMetric(gauge.Name()).SetValue(gauge.Value()))
	}

	return batch
}

func postMetrics(ctx context.Context, client *resty.Client, gz *gzip.Writer, batch []metrics.Metrics) error {
	if len(batch) == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	jsonBody, err := json.Marshal(batch)
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "encode json"),
			logger.Any("data", batch))
		return err
	}

	var hashString string
	if Config.SecretKey() != "" {
		hash := secure.Hash(jsonBody, []byte(Config.SecretKey()))
		hashString = hex.EncodeToString(hash[:])
	}

	if len(publicKey) > 0 {
		jsonBody, err = cryptokey.Encrypt(jsonBody, publicKey)
		if err != nil {
			return err
		}
	}

	buf := bytes.NewBuffer(nil)
	gz.Reset(buf)
	_, err = gz.Write(jsonBody)
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "compress request"),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}
	err = gz.Close()
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "compress request"),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}

	logger.Log.Debug(`Send POST /updates/ request`,
		logger.String("event", "send request"),
		logger.String("addr", Config.ServerAddr()),
		logger.Any("body", json.RawMessage(jsonBody)))

	req := client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(buf)

	if hashString != "" {
		req.SetHeader("HashSHA256", hashString)
	}

	resp, err := req.Post("updates/")

	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "send request"),
			logger.String("url", "http://"+Config.ServerAddr()+"/updates/"),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}

	rawBody := resp.RawBody()
	defer rawBody.Close()

	gzBody, err := gzip.NewReader(rawBody)
	if err != nil {
		return nil
	}
	defer gzBody.Close()

	body, err := io.ReadAll(gzBody)
	if err != nil {
		return nil
	}

	logger.Log.Debug(`Received response from the server`, logger.String("event", "response received"), logger.Any("headers", resp.Header()), logger.Any("body", json.RawMessage(body)))

	return nil
}

func retryPostMetrics(ctx context.Context, client *resty.Client, gz *gzip.Writer, data *storage.MemStorage) error {
	var err error
	var neterr *net.OpError

	// Delay after unsuccessful request
	retryDelay := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
		0,
	}

	for _, delay := range retryDelay {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		batch := packMetricsIntoBatch(data)
		err = postMetrics(ctx, client, gz, batch)

		errors.As(err, &neterr)
		if err == nil || !errors.Is(err, neterr) {
			return nil
		}
		time.Sleep(delay)
	}

	return err
}
