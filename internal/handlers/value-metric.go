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

// ValueMetricHandler processes the request GET /value/{metricType}/{metricID}.
// Returns the metric value.
func ValueMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricType, metricID string

	metricType = chi.URLParam(r, "metricType")
	metricID = chi.URLParam(r, "metricID")

	// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
	if metricType == "" {
		http.Error(w, `Metric type not specified`, http.StatusBadRequest)
		return
	}

	// При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if metricID == "" {
		http.Error(w, `Metric name not specified`, http.StatusNotFound)
		return
	}

	switch metricType {
	case metrics.TypeCounter:
		Counter, ok := storage.CounterContext(r.Context(), metricID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Counter '%s' not found`, metricID), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, strconv.FormatInt(Counter.Value(), 10))
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "value metric handler"), logger.Int64("value", Counter.Value()))
		}

	case metrics.TypeGauge:
		Gauge, ok := storage.GaugeContext(r.Context(), metricID)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Gauge '%s' not found`, metricID), http.StatusNotFound)
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
	}
}
