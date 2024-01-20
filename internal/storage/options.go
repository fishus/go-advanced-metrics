package storage

import "github.com/fishus/go-advanced-metrics/internal/metrics"

type StorageOptions struct {
	gauges   []metrics.Gauge
	counters []metrics.Counter
}

type StorageOption func(o *StorageOptions)

func WithCounters(counters []metrics.Counter) StorageOption {
	return func(o *StorageOptions) {
		o.counters = counters
	}
}

func WithCounter(counter metrics.Counter) StorageOption {
	return func(o *StorageOptions) {
		o.counters = append(o.counters, counter)
	}
}

func WithGauges(gauges []metrics.Gauge) StorageOption {
	return func(o *StorageOptions) {
		o.gauges = gauges
	}
}

func WithGauge(gauge metrics.Gauge) StorageOption {
	return func(o *StorageOptions) {
		o.gauges = append(o.gauges, gauge)
	}
}
