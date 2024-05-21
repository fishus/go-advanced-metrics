package agent

import (
	"context"
	"sync"
	"time"

	cg "github.com/fishus/go-advanced-metrics/internal/agent/client/grpc"
	"github.com/fishus/go-advanced-metrics/internal/agent/client/rest"
	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
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

	var client IAgentClient

	switch Config.ClientType() {
	case ClientTypeREST:
		client = rest.NewClient(rest.Config{
			ServerAddr: Config.ServerAddr(),
			SecretKey:  Config.SecretKey(),
			PublicKey:  PublicKey,
		})
		err := client.Init()
		if err != nil {
			logger.Log.Panic(err.Error())
		}
	case ClientTypeGRPC:
		client = cg.NewClient(cg.Config{
			ServerAddr: Config.ServerAddr(),
			SecretKey:  Config.SecretKey(),
			PublicKey:  PublicKey,
		})
		err := client.Init()
		if err != nil {
			logger.Log.Panic(err.Error())
		}
		conn := client.(*cg.Client).Conn()
		if conn != nil {
			defer conn.Close()
		}
	default:
		logger.Log.Panic("unspecified client type")
	}

	for data := range dataCh {
		batch := packMetricsIntoBatch(data)
		err := client.RetryUpdateBatch(ctx, batch)
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
