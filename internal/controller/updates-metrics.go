package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func (c Controller) UpdatesMetrics(ctx context.Context, metricsBatch []metrics.Metrics) (mb []metrics.Metrics, code int, err error) {
	gaugesBatch := make([]metrics.Gauge, 0)
	countersBatch := make([]metrics.Counter, 0)

	{
		for _, metric := range metricsBatch {
			if err := ValidateInputMetric(metric); err != nil {
				var ve *ValidMetricError
				if errors.As(err, &ve) {
					return mb, ve.HTTPCode, ve
				} else {
					return mb, http.StatusInternalServerError, err
				}
			}

			switch metric.MType {
			case metrics.TypeCounter:
				c, err := metrics.NewCounter(metric.ID, *metric.Delta)
				if err != nil {
					return mb, http.StatusBadRequest, err
				}
				countersBatch = append(countersBatch, *c)
			case metrics.TypeGauge:
				g, err := metrics.NewGauge(metric.ID, *metric.Value)
				if err != nil {
					return mb, http.StatusBadRequest, err
				}
				gaugesBatch = append(gaugesBatch, *g)
			}
		}
	}

	err = Storage.InsertBatchContext(ctx, store.WithCounters(countersBatch), store.WithGauges(gaugesBatch))
	if err != nil {
		return mb, http.StatusInternalServerError, err
	}

	if names := getBatchCounterNames(countersBatch); len(names) > 0 {
		counters := Storage.CountersContext(ctx, store.FilterNames(names))
		for _, cn := range names {
			if c, ok := counters[cn]; ok {
				mb = append(mb, metrics.NewCounterMetric(c.Name()).SetDelta(c.Value()))
			}
		}
	}

	if names := getBatchGaugeNames(gaugesBatch); len(names) > 0 {
		gauges := Storage.GaugesContext(ctx, store.FilterNames(names))
		for _, cn := range names {
			if g, ok := gauges[cn]; ok {
				mb = append(mb, metrics.NewGaugeMetric(g.Name()).SetValue(g.Value()))
			}
		}
	}

	// Synchronously save metrics values into a file
	if s, ok := Storage.(store.SyncSaver); ok {
		err := s.SyncSave()
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "synchronously save metrics into file"))
		}
	}

	return mb, http.StatusOK, nil
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
