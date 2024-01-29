package handlers

import (
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type ValidMetricError struct {
	HTTPCode int
	Err      error
}

func (ve *ValidMetricError) Error() string {
	return ve.Err.Error()
}

func (ve *ValidMetricError) Unwrap() error {
	return ve.Err
}

func NewValidMetricError(httpCode int, err error) *ValidMetricError {
	return &ValidMetricError{
		HTTPCode: httpCode,
		Err:      err,
	}
}

// Функции проверки входящих данных метрики
func validateInputMetric(metric metrics.Metrics) error {
	// При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
	if metric.ID == "" {
		return NewValidMetricError(http.StatusNotFound, errors.New(`ID not specified`))
	}

	// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
	if metric.MType == "" {
		return NewValidMetricError(http.StatusBadRequest, errors.New(`type not specified`))
	}

	switch metric.MType {
	case metrics.TypeCounter:
		if metric.Delta == nil {
			return NewValidMetricError(http.StatusBadRequest, errors.New(`incorrect counter delta`))
		}
	case metrics.TypeGauge:
		if metric.Value == nil {
			return NewValidMetricError(http.StatusBadRequest, errors.New(`incorrect gauge value`))
		}
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		return NewValidMetricError(http.StatusBadRequest, errors.New(`incorrect metric type`))
	}

	return nil
}
