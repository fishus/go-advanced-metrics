package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

// UpdateMetricHandler processes a request like POST /update/{metricType}/{metricID}/{metricValue}
// Stores metric data by type and name
func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics

	metric.ID = chi.URLParam(r, "metricID")
	metric.MType = chi.URLParam(r, "metricType")

	switch metric.MType {
	case metrics.TypeCounter:
		metric.Value = nil

		v := chi.URLParam(r, "metricValue")
		if v == "" {
			metric.Delta = nil
		} else if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			metric.Delta = nil
		} else {
			metric.Delta = new(int64)
			*metric.Delta = i
		}
	case metrics.TypeGauge:
		metric.Delta = nil

		v := chi.URLParam(r, "metricValue")
		if v == "" {
			metric.Value = nil
		} else if f, err := strconv.ParseFloat(v, 64); err != nil {
			metric.Value = nil
		} else {
			metric.Value = new(float64)
			*metric.Value = f
		}
	}

	if err := validateInputMetric(metric); err != nil {
		var ve *ValidMetricError
		if errors.As(err, &ve) {
			JSONError(w, ve.Error(), ve.HTTPCode)
			logger.Log.Debug(ve.Error(), logger.Any("metric", metric))
		} else {
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	switch metric.MType {
	case metrics.TypeCounter:
		err := storage.AddCounterContext(r.Context(), metric.ID, *metric.Delta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
	case metrics.TypeGauge:
		err := storage.SetGaugeContext(r.Context(), metric.ID, *metric.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
	}

	// Synchronously save metrics values into a file
	if s, ok := storage.(store.SyncSaver); ok {
		err := s.SyncSave()
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "synchronously save metrics into file"))
		}
	}

	// При успешном приёме возвращать http.StatusOK.
	w.WriteHeader(http.StatusOK)
}
