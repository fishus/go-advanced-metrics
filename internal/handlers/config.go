package handlers

import "github.com/fishus/go-advanced-metrics/internal/metrics"

type config struct {
	IsSyncMetricsSave bool
}

var Config = config{}

var storage = metrics.NewMemStorage()

func Storage() *metrics.MemStorage {
	return storage
}
