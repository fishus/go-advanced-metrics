package metrics

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) SetValue(value float64) Metrics {
	m.Delta = nil
	if m.Value == nil {
		m.Value = new(float64)
	}
	*m.Value = value
	return m
}

func (m Metrics) SetDelta(delta int64) Metrics {
	m.Value = nil
	if m.Delta == nil {
		m.Delta = new(int64)
	}
	*m.Delta = delta
	return m
}

func NewGaugeMetric(id string) Metrics {
	m := Metrics{
		ID:    id,
		MType: TypeGauge,
		Value: nil,
		Delta: nil,
	}
	return m
}

func NewCounterMetric(id string) Metrics {
	m := Metrics{
		ID:    id,
		MType: TypeCounter,
		Delta: nil,
		Value: nil,
	}
	return m
}
