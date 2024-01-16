package storage

import (
	"encoding/json"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// MemStorage contains a set of values for all metrics and store its in memory
type MemStorage struct {
	gauges   map[string]metrics.Gauge
	counters map[string]metrics.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metrics.Gauge),
		counters: make(map[string]metrics.Counter),
	}
}

// Gauge returns the gauge metric by name
func (m *MemStorage) Gauge(name string) (metrics.Gauge, bool) {
	if v, ok := m.gauges[name]; ok {
		return v, ok
	} else {
		return metrics.Gauge{}, false
	}
}

// GaugeValue returns the gauge metric value by name
func (m *MemStorage) GaugeValue(name string) (float64, bool) {
	if v, ok := m.gauges[name]; ok {
		return v.Value(), ok
	}
	return 0, false
}

// Gauges returns all gauge metrics
func (m *MemStorage) Gauges() map[string]metrics.Gauge {
	return m.gauges
}

func (m *MemStorage) SetGauge(name string, value float64) error {
	if m.gauges == nil {
		m.gauges = make(map[string]metrics.Gauge)
	}
	gauge, ok := m.gauges[name]
	if !ok {
		g, err := metrics.NewGauge(name, value)
		if err != nil {
			return err
		}
		gauge = *g
	} else {
		err := gauge.SetValue(value)
		if err != nil {
			return err
		}
	}
	m.gauges[name] = gauge
	return nil
}

// Counter returns the counter metric by name
func (m *MemStorage) Counter(name string) (metrics.Counter, bool) {
	if v, ok := m.counters[name]; ok {
		return v, ok
	} else {
		return metrics.Counter{}, false
	}
}

// CounterValue returns the counter metric value by name
func (m *MemStorage) CounterValue(name string) (int64, bool) {
	if v, ok := m.counters[name]; ok {
		return v.Value(), ok
	}
	return 0, false
}

// Counters returns all counter metrics
func (m *MemStorage) Counters() map[string]metrics.Counter {
	return m.counters
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	if m.counters == nil {
		m.counters = make(map[string]metrics.Counter)
	}
	counter, ok := m.counters[name]
	if !ok {
		c, err := metrics.NewCounter(name, value)
		if err != nil {
			return err
		}
		counter = *c
	} else {
		err := counter.AddValue(value)
		if err != nil {
			return err
		}
	}
	m.counters[name] = counter
	return nil
}

func (m *MemStorage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Gauges   map[string]metrics.Gauge   `json:"gauges"`
		Counters map[string]metrics.Counter `json:"counters"`
	}{
		Gauges:   m.gauges,
		Counters: m.counters,
	})
}

func (m *MemStorage) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Gauges   map[string]metrics.Gauge   `json:"gauges"`
		Counters map[string]metrics.Counter `json:"counters"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	m.gauges = aux.Gauges
	m.counters = aux.Counters

	return nil
}

var _ MetricsStorager = (*MemStorage)(nil)
var _ json.Marshaler = (*MemStorage)(nil)
var _ json.Unmarshaler = (*MemStorage)(nil)
