package metrics

import "errors"

// Metric contains the values of a single metric
type Metric struct {
	gauge   float64 // Новое значение должно замещать предыдущее
	counter int64   // Новое значение должно добавляться к предыдущему
}

func (m *Metric) SetGauge(value float64) error {
	m.gauge = value
	return nil
}

func (m *Metric) AddCounter(value int64) error {
	if value < 0 {
		return errors.New(`metrics: the counter value must be positive`)
	}
	m.counter += value
	return nil
}

func NewMetric() Metric {
	return Metric{
		gauge:   0.0,
		counter: 0,
	}
}
