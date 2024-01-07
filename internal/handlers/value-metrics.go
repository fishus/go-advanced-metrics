package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// ValueMetricsHandler returns metrics data in JSON format
func ValueMetricsHandler(w http.ResponseWriter, r *http.Request) {
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
		counterValue, ok := storage.CounterValue(metric.ID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			JSONError(w, fmt.Sprintf(`Counter '%s' not found`, metric.ID), http.StatusNotFound)
			logger.Log.Debug(fmt.Sprintf(`Counter '%s' not found`, metric.ID), zap.Any("metric", metric))
			return
		}
		metric.Delta = new(int64)
		*metric.Delta = counterValue
		metric.Value = nil
	case metrics.TypeGauge:
		gaugeValue, ok := storage.GaugeValue(metric.ID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			JSONError(w, fmt.Sprintf(`Gauge '%s' not found`, metric.ID), http.StatusNotFound)
			logger.Log.Debug(fmt.Sprintf(`Gauge '%s' not found`, metric.ID), zap.Any("metric", metric))
			return
		}
		metric.Value = new(float64)
		*metric.Value = gaugeValue
		metric.Delta = nil
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		JSONError(w, `Incorrect metric type`, http.StatusBadRequest)
		logger.Log.Debug(`Incorrect metric type`, zap.String("type", metric.MType))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		logger.Log.Debug(err.Error(), zap.Any("data", metric))
		return
	}
}
