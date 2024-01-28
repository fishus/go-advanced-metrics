package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

// UpdatesMetricsHandler processes a request like POST /updates/
// Store a batch of metrics in JSON format
func UpdatesMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var metricsBatch []metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metricsBatch); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	gaugesBatch := make([]metrics.Gauge, 0)
	countersBatch := make([]metrics.Counter, 0)

	{
		for _, metric := range metricsBatch {
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
				c, err := metrics.NewCounter(metric.ID, *metric.Delta)
				if err != nil {
					JSONError(w, err.Error(), http.StatusBadRequest)
					logger.Log.Debug(err.Error(), logger.Any("metric", metric))
					return
				}
				countersBatch = append(countersBatch, *c)
			case metrics.TypeGauge:
				g, err := metrics.NewGauge(metric.ID, *metric.Value)
				if err != nil {
					JSONError(w, err.Error(), http.StatusBadRequest)
					logger.Log.Debug(err.Error(), logger.Any("metric", metric))
					return
				}
				gaugesBatch = append(gaugesBatch, *g)
			}
		}
	}

	err := storage.InsertBatchContext(r.Context(), store.WithCounters(countersBatch), store.WithGauges(gaugesBatch))
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Debug(err.Error(),
			logger.Any("gaugesBatch", gaugesBatch),
			logger.Any("countersBatch", countersBatch))
		return
	}

	metricsBatch = metricsBatch[:0]

	if names := getBatchCounterNames(countersBatch); len(names) > 0 {
		counters := storage.CountersContext(r.Context(), store.FilterNames(names))
		for _, cn := range names {
			if c, ok := counters[cn]; ok {
				metricsBatch = append(metricsBatch, metrics.NewCounterMetric(c.Name()).SetDelta(c.Value()))
			}
		}
	}

	if names := getBatchGaugeNames(gaugesBatch); len(names) > 0 {
		gauges := storage.GaugesContext(r.Context(), store.FilterNames(names))
		for _, cn := range names {
			if g, ok := gauges[cn]; ok {
				metricsBatch = append(metricsBatch, metrics.NewGaugeMetric(g.Name()).SetValue(g.Value()))
			}
		}
	}

	// Synchronously save metrics values into a file
	if s, ok := storage.(store.SyncSaver); ok {
		err := s.SyncSave()
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "synchronously save metrics into file"))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metricsBatch); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", metricsBatch))
		return
	}
}

func getBatchGaugeNames(batch []metrics.Gauge) []string {
	names := make([]string, 0)
	keys := map[string]bool{}

	for _, m := range batch {
		if !keys[m.Name()] {
			keys[m.Name()] = true
			names = append(names, m.Name())
		}
	}

	return names
}

func getBatchCounterNames(batch []metrics.Counter) []string {
	names := make([]string, 0)
	keys := map[string]bool{}

	for _, m := range batch {
		if !keys[m.Name()] {
			keys[m.Name()] = true
			names = append(names, m.Name())
		}
	}

	return names
}
