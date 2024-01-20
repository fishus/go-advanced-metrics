package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

// UpdatesMetricsHandler processes a request like POST /updates/
// Store a batch of metrics in JSON format
func UpdatesMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type Metric struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}

	var metricsBatch []Metric

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
			if metric.ID == "" {
				JSONError(w, `Metric name not specified`, http.StatusNotFound)
				logger.Log.Debug(`Metric name not specified`)
				return
			}

			if metric.MType == "" {
				JSONError(w, `Metric type not specified`, http.StatusBadRequest)
				logger.Log.Debug(`Metric type not specified`)
				return
			}

			switch metric.MType {
			case metrics.TypeCounter:
				if metric.Delta == nil {
					JSONError(w, `Incorrect counter value`, http.StatusBadRequest)
					logger.Log.Debug(`Incorrect counter value`, logger.Any("metric", metric))
					return
				}

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
				if metric.Value == nil {
					JSONError(w, `Incorrect gauge value`, http.StatusBadRequest)
					logger.Log.Debug(`Incorrect gauge value`, logger.Any("metric", metric))
					return
				}

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
			default:
				JSONError(w, `Incorrect metric type`, http.StatusBadRequest)
				logger.Log.Debug(`Incorrect metric type`, logger.String("type", metric.MType))
				return
			}
		}
	}

	err := storage.InsertBatch(store.WithCounters(countersBatch), store.WithGauges(gaugesBatch))
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Debug(err.Error(),
			logger.Any("gaugesBatch", gaugesBatch),
			logger.Any("countersBatch", countersBatch))
		return
	}

	metricsBatch = metricsBatch[:0]

	if len(cNames) > 0 {
		counters := storage.Counters(store.FilterNames(cNames))
		for _, cn := range cNames {
			if c, ok := counters[cn]; ok {
				metric := Metric{
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
		gauges := storage.Gauges(store.FilterNames(gNames))
		for _, gn := range gNames {
			if g, ok := gauges[gn]; ok {
				metric := Metric{
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
	if Config.IsSyncMetricsSave && fmt.Sprintf("%T", storage) == "*store.FileStorage" {
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
