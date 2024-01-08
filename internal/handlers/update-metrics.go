package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// UpdateMetricsHandler processes a request like POST /update/
// Store data sent in JSON format
func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}

	var metric Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error(), zap.Any("body", r.Body))
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
			logger.Log.Debug(`Incorrect counter value`, zap.Any("metric", metric))
			return
		}

		err := storage.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), zap.Any("metric", metric))
			return
		}
		counterValue, _ := storage.CounterValue(metric.ID)
		metric.Delta = new(int64)
		*metric.Delta = counterValue
		metric.Value = nil
	case metrics.TypeGauge:
		if metric.Value == nil {
			JSONError(w, `Incorrect gauge value`, http.StatusBadRequest)
			logger.Log.Debug(`Incorrect gauge value`, zap.Any("metric", metric))
			return
		}

		err := storage.SetGauge(metric.ID, *metric.Value)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), zap.Any("metric", metric))
			return
		}
		gaugeValue, _ := storage.GaugeValue(metric.ID)
		metric.Value = new(float64)
		*metric.Value = gaugeValue
		metric.Delta = nil
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		JSONError(w, `Incorrect metric type`, http.StatusBadRequest)
		logger.Log.Debug(`Incorrect metric type`, zap.String("type", metric.MType))
		return
	}

	// Save metrics values into a file
	if Config.IsSyncMetricsSave {
		err := storage.Save()
		if !errors.Is(err, metrics.ErrEmptyFilename) {
			if err != nil {
				logger.Log.Error(err.Error(), zap.String("event", "save metrics into file"))
			} else {
				logger.Log.Debug("Metric values saved into file", zap.String("event", "save metrics into file"))
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		logger.Log.Debug(err.Error(), zap.Any("data", metric))
		return
	}
}
