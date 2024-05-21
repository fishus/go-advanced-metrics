package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// UpdatesMetricsHandler processes the request POST /updates/.
// Receives a batch of metrics data in JSON format and stores their values.
func UpdatesMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metricsBatch []metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metricsBatch); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	mb, code, err := Controller.UpdatesMetrics(r.Context(), metricsBatch)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metrics", mb))
		JSONError(w, err.Error(), code)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mb); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", mb))
		return
	}
}
