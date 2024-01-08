package metrics

import "errors"

type Gauge struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func NewGauge(name string, v float64) (*Gauge, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	gauge := &Gauge{Name: name}
	if err := gauge.SetValue(v); err != nil {
		return nil, err
	}
	return gauge, nil
}

func (g *Gauge) SetValue(v float64) error {
	g.Value = v
	return nil
}
