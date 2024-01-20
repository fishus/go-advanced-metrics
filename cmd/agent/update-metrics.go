package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/fishus/go-advanced-metrics/internal/storage"
)

func collectAndSendMetrics() {
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
	now := time.Now()

	pollAfter := now.Add(config.pollInterval)
	reportAfter := now.Add(config.reportInterval)

	for {
		now = time.Now()

		// Collect metrics every {options.pollInterval} seconds
		if now.After(pollAfter) {
			pollAfter = now.Add(config.pollInterval)
			collector.CollectMemStats(ms, data)
		}

		// Send metrics to the server every {options.reportInterval} seconds
		if now.After(reportAfter) {
			reportAfter = now.Add(config.reportInterval)

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

			_ = postUpdateMetrics(client, gz, metricsBatch)

			// Reset metrics
			data = storage.NewMemStorage()
		}

		time.Sleep(1 * time.Second)
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
			logger.ByteString("body", jsonBody))
		return err
	}
	err = gz.Close()
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "compress request"),
			logger.ByteString("body", jsonBody))
		return err
	}

	logger.Log.Debug(`Send POST /updates/ request`,
		logger.String("event", "send request"),
		logger.String("addr", config.serverAddr),
		logger.ByteString("body", jsonBody))

	resp, err := client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post("updates/")

	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "send request"),
			logger.String("url", "http://"+config.serverAddr+"/updates/"),
			logger.ByteString("body", jsonBody))
		return err
	}

	logger.Log.Debug(`Received response from the server`, logger.String("event", "response received"), logger.Any("headers", resp.Header()), logger.Any("body", resp.Body()))

	return nil
}
