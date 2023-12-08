package handlers

import (
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var storage metrics.Repositories = metrics.NewMemStorage()

// UpdateHandler processes a request like POST /update/{metricType}/{metricName}/{metricValue}
// Stores metric data by type and name
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		var metricValue int64

		v := chi.URLParam(r, "metricValue")

		if i, err := strconv.ParseInt(v, 10, 64); err != nil || v == "" {
			http.Error(w, `Incorrect metric value`, http.StatusBadRequest)
			return
		} else {
			metricValue = i
		}

		err := storage.AddCounter(metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case metrics.TypeGauge:
		var metricValue float64

		v := chi.URLParam(r, "metricValue")

		if i, err := strconv.ParseFloat(v, 64); err != nil || v == "" {
			http.Error(w, `Incorrect metric value`, http.StatusBadRequest)
			return
		} else {
			metricValue = i
		}

		err := storage.SetGauge(metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Incorrect metric type`, http.StatusBadRequest)
		return
	}

	// При успешном приёме возвращать http.StatusOK.
	w.WriteHeader(http.StatusOK)
}
