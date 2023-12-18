package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

var config Config

func main() {
	config = loadConfig()
	collectAndSendMetrics()
}

func collectAndSendMetrics() {
	// data contains a set of values for all metrics
	data := metrics.NewMemStorage()

	ms := &runtime.MemStats{}

	client := resty.New()

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
	_, err := client.R().
		SetPathParams(map[string]string{
			"metricType":  mtype,
			"metricName":  name,
			"metricValue": value,
		}).
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", config.serverAddr))

	if err != nil {
		panic(err)
	}

	return nil
}
