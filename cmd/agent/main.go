package main

import (
	"errors"
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const serverHost = "localhost:8080"

func main() {
	// mtx contains a set of values for all metrics
	mtx := make(map[string]metrics.Metric)

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

		runtime.ReadMemStats(ms)

		_ = setMetricGauge(mtx, "Alloc", float64(ms.Alloc))
		_ = setMetricGauge(mtx, "BuckHashSys", float64(ms.BuckHashSys))
		_ = setMetricGauge(mtx, "Frees", float64(ms.Frees))
		_ = setMetricGauge(mtx, "GCCPUFraction", float64(ms.GCCPUFraction))
		_ = setMetricGauge(mtx, "GCSys", float64(ms.GCSys))
		_ = setMetricGauge(mtx, "HeapAlloc", float64(ms.HeapAlloc))
		_ = setMetricGauge(mtx, "HeapIdle", float64(ms.HeapIdle))
		_ = setMetricGauge(mtx, "HeapInuse", float64(ms.HeapInuse))
		_ = setMetricGauge(mtx, "HeapObjects", float64(ms.HeapObjects))
		_ = setMetricGauge(mtx, "HeapReleased", float64(ms.HeapReleased))
		_ = setMetricGauge(mtx, "HeapSys", float64(ms.HeapSys))
		_ = setMetricGauge(mtx, "LastGC", float64(ms.LastGC))
		_ = setMetricGauge(mtx, "Lookups", float64(ms.Lookups))
		_ = setMetricGauge(mtx, "MCacheInuse", float64(ms.MCacheInuse))
		_ = setMetricGauge(mtx, "MCacheSys", float64(ms.MCacheSys))
		_ = setMetricGauge(mtx, "MSpanInuse", float64(ms.MSpanInuse))
		_ = setMetricGauge(mtx, "MSpanSys", float64(ms.MSpanSys))
		_ = setMetricGauge(mtx, "Mallocs", float64(ms.Mallocs))
		_ = setMetricGauge(mtx, "NextGC", float64(ms.NextGC))
		_ = setMetricGauge(mtx, "NumForcedGC", float64(ms.NumForcedGC))
		_ = setMetricGauge(mtx, "NumGC", float64(ms.NumGC))
		_ = setMetricGauge(mtx, "OtherSys", float64(ms.OtherSys))
		_ = setMetricGauge(mtx, "PauseTotalNs", float64(ms.PauseTotalNs))
		_ = setMetricGauge(mtx, "StackInuse", float64(ms.StackInuse))
		_ = setMetricGauge(mtx, "StackSys", float64(ms.StackSys))
		_ = setMetricGauge(mtx, "Sys", float64(ms.Sys))
		_ = setMetricGauge(mtx, "TotalAlloc", float64(ms.TotalAlloc))

		_ = addMetricCounter(mtx, "PollCount", 1)
		_ = setMetricGauge(mtx, "RandomValue", rand.Float64()*100)

		if i%reportInterval != 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for k, v := range mtx {
			g := v.Gauge()
			c := v.Counter()

			if g != 0 {
				_ = postUpdateMetrics(client, "gauge", k, strconv.FormatFloat(g, 'f', -1, 64))
			}

			if c != 0 {
				_ = postUpdateMetrics(client, "counter", k, strconv.FormatInt(c, 10))
			}
		}

		// Reset counters and data
		i = 0
		mtx = make(map[string]metrics.Metric)

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

func setMetricGauge(mtx map[string]metrics.Metric, name string, value float64) error {
	if mtx == nil {
		return errors.New("mtx map[string]metrics.Metric is nil")
	}
	m, ok := mtx[name]
	if !ok {
		m = metrics.NewMetric()
	}

	err := m.SetGauge(value)
	if err != nil {
		return err
	}

	mtx[name] = m
	return nil
}

func addMetricCounter(mtx map[string]metrics.Metric, name string, value int64) error {
	if mtx == nil {
		return errors.New("mtx map[string]metrics.Metric is nil")
	}
	m, ok := mtx[name]
	if !ok {
		m = metrics.NewMetric()
	}

	err := m.AddCounter(value)
	if err != nil {
		return err
	}

	mtx[name] = m
	return nil
}
