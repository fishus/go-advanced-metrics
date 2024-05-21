package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// UpdateMetricsHandler processes the request POST /update/.
// Receives metric data and stores its value.
func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	m, code, err := Controller.UpdateMetrics(r.Context(), metric)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metric", metric))
		JSONError(w, err.Error(), code)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(m); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", m))
		return
	}
}
