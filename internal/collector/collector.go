package collector

import (
	"context"
	"errors"
	"math/rand"
	"runtime"

	"github.com/fishus/go-advanced-metrics/internal/storage"
)

func CollectMemStats(ctx context.Context) *storage.MemStorage {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	ms := &runtime.MemStats{}
	data := storage.NewMemStorage()

	runtime.ReadMemStats(ms)

	_ = setMetricGauge(data, "Alloc", float64(ms.Alloc))
	_ = setMetricGauge(data, "BuckHashSys", float64(ms.BuckHashSys))
	_ = setMetricGauge(data, "Frees", float64(ms.Frees))
	_ = setMetricGauge(data, "GCCPUFraction", float64(ms.GCCPUFraction))
	_ = setMetricGauge(data, "GCSys", float64(ms.GCSys))
	_ = setMetricGauge(data, "HeapAlloc", float64(ms.HeapAlloc))
	_ = setMetricGauge(data, "HeapIdle", float64(ms.HeapIdle))
	_ = setMetricGauge(data, "HeapInuse", float64(ms.HeapInuse))
	_ = setMetricGauge(data, "HeapObjects", float64(ms.HeapObjects))
	_ = setMetricGauge(data, "HeapReleased", float64(ms.HeapReleased))
	_ = setMetricGauge(data, "HeapSys", float64(ms.HeapSys))
	_ = setMetricGauge(data, "LastGC", float64(ms.LastGC))
	_ = setMetricGauge(data, "Lookups", float64(ms.Lookups))
	_ = setMetricGauge(data, "MCacheInuse", float64(ms.MCacheInuse))
	_ = setMetricGauge(data, "MCacheSys", float64(ms.MCacheSys))
	_ = setMetricGauge(data, "MSpanInuse", float64(ms.MSpanInuse))
	_ = setMetricGauge(data, "MSpanSys", float64(ms.MSpanSys))
	_ = setMetricGauge(data, "Mallocs", float64(ms.Mallocs))
	_ = setMetricGauge(data, "NextGC", float64(ms.NextGC))
	_ = setMetricGauge(data, "NumForcedGC", float64(ms.NumForcedGC))
	_ = setMetricGauge(data, "NumGC", float64(ms.NumGC))
	_ = setMetricGauge(data, "OtherSys", float64(ms.OtherSys))
	_ = setMetricGauge(data, "PauseTotalNs", float64(ms.PauseTotalNs))
	_ = setMetricGauge(data, "StackInuse", float64(ms.StackInuse))
	_ = setMetricGauge(data, "StackSys", float64(ms.StackSys))
	_ = setMetricGauge(data, "Sys", float64(ms.Sys))
	_ = setMetricGauge(data, "TotalAlloc", float64(ms.TotalAlloc))

	_ = addMetricCounter(data, "PollCount", 1)
	_ = setMetricGauge(data, "RandomValue", rand.Float64()*100)

	return data
}

func setMetricGauge(data *storage.MemStorage, name string, value float64) error {
	if data == nil {
		return errors.New("data is nil")
	}

	err := data.SetGauge(name, value)
	if err != nil {
		return err
	}
	return nil
}

func addMetricCounter(data *storage.MemStorage, name string, value int64) error {
	if data == nil {
		return errors.New("data is nil")
	}

	err := data.AddCounter(name, value)
	if err != nil {
		return err
	}
	return nil
}
