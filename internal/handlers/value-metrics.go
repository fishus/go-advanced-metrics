package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// ValueMetricsHandler returns metrics data in JSON format
func ValueMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var metric metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
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
		counterValue, ok := storage.CounterValueContext(r.Context(), metric.ID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			JSONError(w, fmt.Sprintf(`Counter '%s' not found`, metric.ID), http.StatusNotFound)
			logger.Log.Debug(fmt.Sprintf(`Counter '%s' not found`, metric.ID), logger.Any("metric", metric))
			return
		}
		metric = metric.SetDelta(counterValue)
	case metrics.TypeGauge:
		gaugeValue, ok := storage.GaugeValueContext(r.Context(), metric.ID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			JSONError(w, fmt.Sprintf(`Gauge '%s' not found`, metric.ID), http.StatusNotFound)
			logger.Log.Debug(fmt.Sprintf(`Gauge '%s' not found`, metric.ID), logger.Any("metric", metric))
			return
		}
		metric = metric.SetValue(gaugeValue)
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		JSONError(w, `Incorrect metric type`, http.StatusBadRequest)
		logger.Log.Debug(`Incorrect metric type`, logger.String("type", metric.MType))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", metric))
		return
	}
}
