package metrics

import (
	"encoding/json"
	"errors"
)

type Gauge struct {
	name  string
	value float64
}

func NewGauge(name string, v float64) (*Gauge, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	gauge := Gauge{name: name}
	err := gauge.SetValue(v)
	if err != nil {
		return nil, err
	}
	return &gauge, nil
}

func (g Gauge) Name() string {
	return g.name
}

func (g Gauge) Value() float64 {
	return g.value
}

func (g *Gauge) SetValue(v float64) error {
	g.value = v
	return nil
}

func (g Gauge) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}{
		Name:  g.name,
		Value: g.value,
	})
}

func (g *Gauge) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	g.name = aux.Name
	g.value = aux.Value

	return nil
}
