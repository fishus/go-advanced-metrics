package storage

import "github.com/fishus/go-advanced-metrics/internal/metrics"

type GaugeStorager interface {
	Gauge(name string) (metrics.Gauge, bool)
	GaugeValue(name string) (float64, bool)
	Gauges() map[string]metrics.Gauge
	SetGauge(name string, value float64) error
}

type CounterStorager interface {
	Counter(name string) (metrics.Counter, bool)
	CounterValue(name string) (int64, bool)
	Counters() map[string]metrics.Counter
	AddCounter(name string, value int64) error
}

// MetricsStorager is an interface for managing a set of metrics
type MetricsStorager interface {
	GaugeStorager
	CounterStorager
}

// StorageSaver is an interface for save a set of metrics
type StorageSaver interface {
	Save() error
}

// StorageLoader is an interface for load a set of metrics
type StorageLoader interface {
	Load() error
}
