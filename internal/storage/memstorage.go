package storage

import (
	"context"
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
	if gauge, ok := m.gauges[name]; ok {
		return gauge.Value(), ok
	}
	return 0, false
}

// Gauges returns all gauge metrics
func (m *MemStorage) Gauges(filters ...StorageFilter) map[string]metrics.Gauge {
	f := &StorageFilters{}
	for _, filter := range filters {
		filter(f)
	}

	if len(f.names) > 0 {
		diff := make(map[string]metrics.Gauge)

		for _, name := range f.names {
			if g, ok := m.gauges[name]; ok {
				diff[name] = g
			}
		}

		return diff
	}

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
func (m *MemStorage) Counters(filters ...StorageFilter) map[string]metrics.Counter {
	f := &StorageFilters{}
	for _, filter := range filters {
		filter(f)
	}

	if len(f.names) > 0 {
		diff := make(map[string]metrics.Counter)

		for _, name := range f.names {
			if c, ok := m.counters[name]; ok {
				diff[name] = c
			}
		}

		return diff
	}

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

func (m *MemStorage) InsertBatch(opts ...StorageOption) error {
	return m.InsertBatchContext(context.Background(), opts...)
}

func (m *MemStorage) InsertBatchContext(ctx context.Context, opts ...StorageOption) error {
	o := &StorageOptions{}
	for _, opt := range opts {
		opt(o)
	}

	if len(o.gauges) == 0 && len(o.counters) == 0 {
		return nil
	}

	// Сначала пробуем добавить во временное пустое хранилище.
	// Если будут ошибки - откат. Без ошибок сохраняем в настоящее хранилище
	{
		// Temp storage
		ts := &MemStorage{}

		// Check counters for errors
		if len(o.counters) > 0 {
			for _, c := range o.counters {
				err := ts.AddCounter(c.Name(), c.Value())
				if err != nil {
					return err
				}
			}
		}

		// Check gauges for errors
		if len(o.gauges) > 0 {
			for _, g := range o.gauges {
				err := ts.SetGauge(g.Name(), g.Value())
				if err != nil {
					return err
				}
			}
		}
	}

	// Insert to storage
	{
		if len(o.counters) > 0 {
			for _, c := range o.counters {
				_ = m.AddCounter(c.Name(), c.Value())
			}
		}

		if len(o.gauges) > 0 {
			for _, g := range o.gauges {
				_ = m.SetGauge(g.Name(), g.Value())
			}
		}
	}

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
