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

// UpdatesMetricsHandler processes a request like POST /updates/
// Store a batch of metrics in JSON format
func UpdatesMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var metricsBatch []metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metricsBatch); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error(), logger.Any("body", r.Body))
		return
	}

	gaugesBatch := make([]metrics.Gauge, 0)
	countersBatch := make([]metrics.Counter, 0)

	gNames := make([]string, 0)
	cNames := make([]string, 0)

	{
		gMap := map[string]bool{}
		cMap := map[string]bool{}

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
				if !cMap[c.Name()] {
					cMap[c.Name()] = true
					cNames = append(cNames, c.Name())
				}
			case metrics.TypeGauge:
				g, err := metrics.NewGauge(metric.ID, *metric.Value)
				if err != nil {
					JSONError(w, err.Error(), http.StatusBadRequest)
					logger.Log.Debug(err.Error(), logger.Any("metric", metric))
					return
				}
				gaugesBatch = append(gaugesBatch, *g)
				if !gMap[g.Name()] {
					gMap[g.Name()] = true
					gNames = append(gNames, g.Name())
				}
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

	if len(cNames) > 0 {
		counters := storage.CountersContext(r.Context(), store.FilterNames(cNames))
		for _, cn := range cNames {
			if c, ok := counters[cn]; ok {
				metric := metrics.Metrics{
					ID:    c.Name(),
					MType: metrics.TypeCounter,
					Delta: new(int64),
				}
				*metric.Delta = c.Value()
				metricsBatch = append(metricsBatch, metric)
			}
		}
	}

	if len(gNames) > 0 {
		gauges := storage.GaugesContext(r.Context(), store.FilterNames(gNames))
		for _, gn := range gNames {
			if g, ok := gauges[gn]; ok {
				metric := metrics.Metrics{
					ID:    g.Name(),
					MType: metrics.TypeGauge,
					Value: new(float64),
				}
				*metric.Value = g.Value()
				metricsBatch = append(metricsBatch, metric)
			}
		}
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

	if err := json.NewEncoder(w).Encode(metricsBatch); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", metricsBatch))
		return
	}
}
