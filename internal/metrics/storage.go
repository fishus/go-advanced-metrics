package metrics

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrEmptyFilename = errors.New("filename for store metrics data is empty")

// MemStorage contains a set of values for all metrics
type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
	Filename string
	muSave   sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]Gauge),
		counters: make(map[string]Counter),
	}
}

// Save saves metric values to a file.
func (m *MemStorage) Save() error {
	m.muSave.Lock()
	defer m.muSave.Unlock()

	if m.Filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.OpenFile(m.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(&m)
	if err != nil {
		return err
	}

	return nil
}

// Load reads metric values from a file.
func (m *MemStorage) Load() error {
	if m.Filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.OpenFile(m.Filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&m); err != nil {
		return err
	}

	return nil
}

func (m *MemStorage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Gauges   map[string]Gauge   `json:"gauges"`
		Counters map[string]Counter `json:"counters"`
	}{
		Gauges:   m.gauges,
		Counters: m.counters,
	})
}

func (m *MemStorage) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Gauges   map[string]Gauge   `json:"gauges"`
		Counters map[string]Counter `json:"counters"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	m.gauges = aux.Gauges
	m.counters = aux.Counters

	return nil
}

type GaugeRepositories interface {
	Gauge(name string) (Gauge, bool)
	GaugeValue(name string) (float64, bool)
	Gauges() map[string]Gauge
	SetGauge(name string, value float64) error
}

var _ GaugeRepositories = (*MemStorage)(nil)

type CounterRepositories interface {
	Counter(name string) (Counter, bool)
	CounterValue(name string) (int64, bool)
	Counters() map[string]Counter
	AddCounter(name string, value int64) error
}

var _ CounterRepositories = (*MemStorage)(nil)

// Repositories is an interface for managing a set of metrics
type Repositories interface {
	GaugeRepositories
	CounterRepositories
	Save() error
	Load() error
}

var _ Repositories = (*MemStorage)(nil)

// Gauge returns the gauge metric by name
func (m *MemStorage) Gauge(name string) (Gauge, bool) {
	if v, ok := m.gauges[name]; ok {
		return v, ok
	} else {
		return Gauge{}, false
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
func (m *MemStorage) Gauges() map[string]Gauge {
	return m.gauges
}

func (m *MemStorage) SetGauge(name string, value float64) error {
	if m.gauges == nil {
		m.gauges = make(map[string]Gauge)
	}
	gauge, ok := m.gauges[name]
	if !ok {
		g, err := NewGauge(name, value)
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
func (m *MemStorage) Counter(name string) (Counter, bool) {
	if v, ok := m.counters[name]; ok {
		return v, ok
	} else {
		return Counter{}, false
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
func (m *MemStorage) Counters() map[string]Counter {
	return m.counters
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	if m.counters == nil {
		m.counters = make(map[string]Counter)
	}
	counter, ok := m.counters[name]
	if !ok {
		c, err := NewCounter(name, value)
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
