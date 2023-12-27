package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

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
		metricValue, ok := storage.Counter(metricName)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Counter '%s' not found`, metricName), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, strconv.FormatInt(int64(metricValue), 10))
		if err != nil {
			logger.Log.Error(err.Error(), zap.String("event", "value metric handler"), zap.Int64("value", int64(metricValue)))
		}

	case metrics.TypeGauge:
		metricValue, ok := storage.Gauge(metricName)
		if !ok {
			// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
			http.Error(w, fmt.Sprintf(`Gauge '%s' not found`, metricName), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, strconv.FormatFloat(float64(metricValue), 'f', -1, 64))
		if err != nil {
			logger.Log.Error(err.Error(), zap.String("event", "value metric handler"), zap.Float64("value", float64(metricValue)))
		}

	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Incorrect metric type`, http.StatusBadRequest)
		return
	}
}
