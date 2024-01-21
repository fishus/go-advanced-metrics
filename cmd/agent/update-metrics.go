package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/fishus/go-advanced-metrics/internal/storage"
)

func collectAndPostMetrics() {
	// data contains a set of values for all metrics
	data := storage.NewMemStorage()

	ms := &runtime.MemStats{}

	client := resty.New().SetBaseURL("http://" + config.serverAddr)
	logger.Log.Info("Running agent", logger.String("address", config.serverAddr), logger.String("event", "start agent"))

	gz, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		logger.Log.Panic(err.Error())
		return
	}

	ctx := context.Background()

	go collectMetricsAtIntervals(ctx, ms, data)
	go postMetricsAtIntervals(ctx, client, gz, data)

	<-ctx.Done()
}

// Collect metrics every {options.pollInterval} seconds
func collectMetricsAtIntervals(ctx context.Context, ms *runtime.MemStats, data *storage.MemStorage) {
	ticker := time.NewTicker(config.pollInterval)
	for {
		<-ticker.C
		collector.CollectMemStats(ms, data)
	}
}

func postMetricsAtIntervals(ctx context.Context, client *resty.Client, gz *gzip.Writer, data *storage.MemStorage) {
	ticker := time.NewTicker(config.reportInterval)
	for {
		<-ticker.C
		metricsBatch := make([]metrics.Metrics, 0)

		for name, counter := range data.Counters() {
			metric := metrics.Metrics{
				ID:    name,
				MType: metrics.TypeCounter,
				Delta: new(int64),
				Value: nil,
			}
			*metric.Delta = counter.Value()
			metricsBatch = append(metricsBatch, metric)
		}

		for name, gauge := range data.Gauges() {
			metric := metrics.Metrics{
				ID:    name,
				MType: metrics.TypeGauge,
				Value: new(float64),
				Delta: nil,
			}
			*metric.Value = gauge.Value()
			metricsBatch = append(metricsBatch, metric)
		}
		err := postUpdateMetrics(client, gz, metricsBatch)
		if err == nil {
			// Reset metrics
			_ = data.Reset()
		}
	}
}

func postUpdateMetrics(client *resty.Client, gz *gzip.Writer, batch []metrics.Metrics) error {
	jsonBody, err := json.Marshal(batch)
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "encode json"),
			logger.Any("data", batch))
		return err
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
		logger.String("addr", config.serverAddr),
		logger.Any("body", json.RawMessage(jsonBody)))

	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(buf).
		Post("updates/")

	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "send request"),
			logger.String("url", "http://"+config.serverAddr+"/updates/"),
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
