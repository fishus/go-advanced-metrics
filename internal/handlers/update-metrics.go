package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

// UpdateMetricsHandler processes a request like POST /update/
// Store data sent in JSON format
func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var metric metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error(), logger.Any("body", r.Body))
		return
	}

	// При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if metric.ID == "" {
		JSONError(w, `Metric name not specified`, http.StatusNotFound)
		logger.Log.Debug(`Metric name not specified`)
		return
	}

	// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
	if metric.MType == "" {
		JSONError(w, `Metric type not specified`, http.StatusBadRequest)
		logger.Log.Debug(`Metric type not specified`)
		return
	}

	switch metric.MType {
	case metrics.TypeCounter:
		if metric.Delta == nil {
			JSONError(w, `Incorrect counter value`, http.StatusBadRequest)
			logger.Log.Debug(`Incorrect counter value`, logger.Any("metric", metric))
			return
		}

		err := storage.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
		counterValue, _ := storage.CounterValue(metric.ID)
		metric.Delta = new(int64)
		*metric.Delta = counterValue
		metric.Value = nil
	case metrics.TypeGauge:
		if metric.Value == nil {
			JSONError(w, `Incorrect gauge value`, http.StatusBadRequest)
			logger.Log.Debug(`Incorrect gauge value`, logger.Any("metric", metric))
			return
		}

		err := storage.SetGauge(metric.ID, *metric.Value)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
		gaugeValue, _ := storage.GaugeValue(metric.ID)
		metric.Value = new(float64)
		*metric.Value = gaugeValue
		metric.Delta = nil
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		JSONError(w, `Incorrect metric type`, http.StatusBadRequest)
		logger.Log.Debug(`Incorrect metric type`, logger.String("type", metric.MType))
		return
	}

	// Save metrics values into a file
	if Config.IsSyncMetricsSave && fmt.Sprintf("%T", storage) == "*store.FileStorage" {
		err := storage.(*store.FileStorage).Save()
		if !errors.Is(err, store.ErrEmptyFilename) {
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
			} else {
				logger.Log.Debug("Metric values saved into file", logger.String("event", "save metrics into file"))
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", metric))
		return
	}
}
