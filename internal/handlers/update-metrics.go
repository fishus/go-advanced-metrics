package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

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
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
		counterValue, _ := storage.CounterValueContext(r.Context(), metric.ID)
		metric.Delta = new(int64)
		*metric.Delta = counterValue
		metric.Value = nil
	case metrics.TypeGauge:
		err := storage.SetGaugeContext(r.Context(), metric.ID, *metric.Value)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return
		}
		gaugeValue, _ := storage.GaugeValueContext(r.Context(), metric.ID)
		metric.Value = new(float64)
		*metric.Value = gaugeValue
		metric.Delta = nil
	}

	// Save metrics values into a file
	if Config.IsSyncMetricsSave && reflect.TypeOf(storage).String() == reflect.TypeOf((*store.FileStorage)(nil)).String() {
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
