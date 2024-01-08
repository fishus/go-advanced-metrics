package handlers

import "github.com/fishus/go-advanced-metrics/internal/metrics"

var storage = metrics.NewMemStorage()

func Storage() *metrics.MemStorage {
	return storage
}
