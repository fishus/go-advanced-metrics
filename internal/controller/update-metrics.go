package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func (c Controller) UpdateMetrics(ctx context.Context, metric metrics.Metrics) (m metrics.Metrics, code int, err error) {
	m = metric
	if err = ValidateInputMetric(metric); err != nil {
		var ve *ValidMetricError
		if errors.As(err, &ve) {
			return m, ve.HTTPCode, ve
		} else {
			return m, http.StatusInternalServerError, err
		}
	}

	switch metric.MType {
	case metrics.TypeCounter:
		err = Storage.AddCounterContext(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return m, http.StatusBadRequest, err
		}
		counterValue, _ := Storage.CounterValueContext(ctx, metric.ID)
		m = metric.SetDelta(counterValue)
	case metrics.TypeGauge:
		err = Storage.SetGaugeContext(ctx, metric.ID, *metric.Value)
		if err != nil {
			return m, http.StatusBadRequest, err
		}
		gaugeValue, _ := Storage.GaugeValueContext(ctx, metric.ID)
		m = metric.SetValue(gaugeValue)
	}

	// Synchronously save metrics values into a file
	if s, ok := Storage.(store.SyncSaver); ok {
		err := s.SyncSave()
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "synchronously save metrics into file"))
		}
	}

	return m, http.StatusOK, nil
}
