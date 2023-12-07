package main

import (
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/collector"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const serverHost = "localhost:8080"

func main() {
	// data contains a set of values for all metrics
	data := metrics.NewMemStorage()

	ms := &runtime.MemStats{}

	client := &http.Client{}

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
			_ = postUpdateMetrics(client, string(metrics.TypeGauge), name, strconv.FormatFloat(float64(g), 'f', -1, 64))
		}

		for name, c := range data.Counters() {
			_ = postUpdateMetrics(client, string(metrics.TypeCounter), name, strconv.FormatInt(int64(c), 10))
		}

		// Reset counters and data
		i = 0
		data = metrics.NewMemStorage()

		time.Sleep(1 * time.Second)
	}
}

func postUpdateMetrics(client *http.Client, mtype, name, value string) error {
	requestURL := fmt.Sprintf("http://%s/update/%s/%s/%s", serverHost, mtype, name, value)

	resp, err := client.Post(requestURL, "text/plain; charset=utf-8", nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
