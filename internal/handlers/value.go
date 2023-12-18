package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

func ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf(`%s method not allowed`, r.Method), http.StatusMethodNotAllowed)
		return
	}

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
			panic(err)
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
			panic(err)
		}

	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Incorrect metric type`, http.StatusBadRequest)
		return
	}
}
