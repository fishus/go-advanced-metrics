package main

import (
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	collectAndSendMetrics()
}

func collectAndSendMetrics() {
	// data contains a set of values for all metrics
	data := metrics.NewMemStorage()

	ms := &runtime.MemStats{}

	client := resty.New()
	logger.Log.Info("Running agent", zap.String("address", config.serverAddr), zap.String("event", "start agent"))

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

			for name, c := range data.Counters() {
				_ = postUpdateMetrics(client, metrics.TypeCounter, name, int64(c), 0)
			}

			for name, g := range data.Gauges() {
				_ = postUpdateMetrics(client, metrics.TypeGauge, name, 0, float64(g))
			}

			// Reset metrics
			data = metrics.NewMemStorage()
		}

		time.Sleep(1 * time.Second)
	}
}

func postUpdateMetrics(client *resty.Client, mtype, name string, delta int64, value float64) error {
	type Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}

	var metric Metrics

	metric.ID = name
	metric.MType = mtype

	switch mtype {
	case metrics.TypeCounter:
		metric.Delta = new(int64)
		*metric.Delta = delta
	case metrics.TypeGauge:
		metric.Value = new(float64)
		*metric.Value = value
	}

	logger.Log.Debug(`Sent POST /update/ request`, zap.String("event", "request sent"), zap.String("addr", config.serverAddr), zap.Any("metric", metric))

	resp, err := client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(metric).
		Post("http://" + config.serverAddr + "/update/")

	if err != nil {
		logger.Log.Error(err.Error(),
			zap.String("url", "http://"+config.serverAddr+"/update/"),
			zap.Any("data", metric))
	}

	logger.Log.Debug(`Received response from the server`, zap.String("event", "response received"), zap.Any("headers", resp.Header()), zap.Any("body", resp.Body()))

	return nil
}
