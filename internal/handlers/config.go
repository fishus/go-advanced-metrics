package handlers

import (
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

type config struct {
	IsSyncMetricsSave bool // Сохранять значения метрик синхронно в файл
}

var Config = config{}

var storage store.MetricsStorager

func Storage() store.MetricsStorager {
	return storage
}

func SetStorage(s store.MetricsStorager) {
	storage = s
}
