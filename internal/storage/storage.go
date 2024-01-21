package storage

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type GaugeStorager interface {
	Gauge(name string) (metrics.Gauge, bool)
	GaugeContext(ctx context.Context, name string) (metrics.Gauge, bool)
	GaugeValue(name string) (float64, bool)
	GaugeValueContext(ctx context.Context, name string) (float64, bool)
	Gauges(filters ...StorageFilter) map[string]metrics.Gauge
	GaugesContext(ctx context.Context, filters ...StorageFilter) map[string]metrics.Gauge
	SetGauge(name string, value float64) error
	SetGaugeContext(ctx context.Context, name string, value float64) error
	ResetGauges() error
}

type CounterStorager interface {
	Counter(name string) (metrics.Counter, bool)
	CounterContext(ctx context.Context, name string) (metrics.Counter, bool)
	CounterValue(name string) (int64, bool)
	CounterValueContext(ctx context.Context, name string) (int64, bool)
	Counters(filters ...StorageFilter) map[string]metrics.Counter
	CountersContext(ctx context.Context, filters ...StorageFilter) map[string]metrics.Counter
	AddCounter(name string, value int64) error
	AddCounterContext(ctx context.Context, name string, value int64) error
	ResetCounters() error
}

// MetricsStorager is an interface for managing a set of metrics
type MetricsStorager interface {
	GaugeStorager
	CounterStorager
	Reset() error
	InsertBatch(opts ...StorageOption) error
	InsertBatchContext(ctx context.Context, opts ...StorageOption) error
}

// StorageSaver is an interface for save a set of metrics
type StorageSaver interface {
	Save() error
}

// StorageLoader is an interface for load a set of metrics
type StorageLoader interface {
	Load() error
}
