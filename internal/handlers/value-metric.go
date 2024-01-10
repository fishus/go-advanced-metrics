package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// ValueMetricHandler returns metrics data
func ValueMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricType, metricName string

	metricType = chi.URLParam(r, "metricType")
	metricName = chi.URLParam(r, "metricName")

	// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
	if metricType == "" {
		http.Error(w, `Metric type not specified`, http.StatusBadRequest)
		return
	}

	// При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if metricName == "" {
		http.Error(w, `Metric name not specified`, http.StatusNotFound)
		return
	}

	switch metricType {
	case metrics.TypeCounter:
		Counter, ok := storage.Counter(metricName)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Counter '%s' not found`, metricName), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, strconv.FormatInt(Counter.Value(), 10))
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "value metric handler"), logger.Int64("value", Counter.Value()))
		}

	case metrics.TypeGauge:
		Gauge, ok := storage.Gauge(metricName)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Gauge '%s' not found`, metricName), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, strconv.FormatFloat(Gauge.Value(), 'f', -1, 64))
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "value metric handler"), logger.Float64("value", Gauge.Value()))
		}

	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Incorrect metric type`, http.StatusBadRequest)
		return
	}
}
