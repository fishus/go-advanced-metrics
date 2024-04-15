package metrics

// Available metric types.
const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

// Metrics structure is used to process incoming data and return results in handlers.
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

// SetValue sets a new value for a metric of type gauge.
func (m Metrics) SetValue(value float64) Metrics {
	m.Delta = nil
	if m.Value == nil {
		m.Value = new(float64)
	}
	*m.Value = value
	return m
}

// SetDelta sets a new value for a metric of type counter.
func (m Metrics) SetDelta(delta int64) Metrics {
	m.Value = nil
	if m.Delta == nil {
		m.Delta = new(int64)
	}
	*m.Delta = delta
	return m
}

// NewGaugeMetric returns a metric data structure of type gauge.
//
// The value of the Metrics.MType field is TypeGauge ("gauge").
// Metric data of type gauge is stored in the Metrics.Value field.
func NewGaugeMetric(id string) Metrics {
	m := Metrics{
		ID:    id,
		MType: TypeGauge,
		Value: nil,
		Delta: nil,
	}
	return m
}

// NewCounterMetric returns a metric data structure of type counter.
//
// The value of the Metrics.MType field is TypeCounter ("counter").
// Metric data of type counter is stored in the Metrics.Delta field.
func NewCounterMetric(id string) Metrics {
	m := Metrics{
		ID:    id,
		MType: TypeCounter,
		Delta: nil,
		Value: nil,
	}
	return m
}
