package main

import (
	"fmt"
	"runtime"
	"strconv"
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

			for name, g := range data.Gauges() {
				_ = postUpdateMetrics(client, metrics.TypeGauge, name, strconv.FormatFloat(float64(g), 'f', -1, 64))
			}

			for name, c := range data.Counters() {
				_ = postUpdateMetrics(client, metrics.TypeCounter, name, strconv.FormatInt(int64(c), 10))
			}

			// Reset metrics
			data = metrics.NewMemStorage()
		}

		time.Sleep(1 * time.Second)
	}
}

func postUpdateMetrics(client *resty.Client, mtype, name, value string) error {
	logger.Log.Debug(`Sent POST /update/ request`, zap.String("event", "request sent"), zap.String("addr", config.serverAddr), zap.String("name", name), zap.String("type", mtype), zap.String("value", value))

	resp, err := client.R().
		SetPathParams(map[string]string{
			"metricType":  mtype,
			"metricName":  name,
			"metricValue": value,
		}).
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", config.serverAddr))

	if err != nil {
		logger.Log.Error(err.Error(),
			zap.String("url", "http://"+config.serverAddr+"/update/"))
	}

	logger.Log.Debug(`Received response from the server`, zap.String("event", "response received"), zap.Any("headers", resp.Header()), zap.Any("body", resp.Body()))

	return nil
}
