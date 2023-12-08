package main

import (
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"runtime"
	"strconv"
	"time"
)

const serverHost = "localhost:8080"

func main() {
	// data contains a set of values for all metrics
	data := metrics.NewMemStorage()

	ms := &runtime.MemStats{}

	client := resty.New()

	// Обновлять метрики с заданной частотой: pollInterval — 2 секунды.
	pollInterval := 2

	// Отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
	reportInterval := 10

	i := 0

	for {
		i++

		if i%pollInterval != 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		collector.CollectMemStats(ms, data)

		if i%reportInterval != 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for name, g := range data.Gauges() {
			_ = postUpdateMetrics(client, metrics.TypeGauge, name, strconv.FormatFloat(float64(g), 'f', -1, 64))
		}

		for name, c := range data.Counters() {
			_ = postUpdateMetrics(client, metrics.TypeCounter, name, strconv.FormatInt(int64(c), 10))
		}

		// Reset counters and data
		i = 0
		data = metrics.NewMemStorage()

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
		Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", serverHost))

	if err != nil {
		panic(err)
	}

	return nil
}
