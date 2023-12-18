package metrics

// MemStorage contains a set of values for all metrics
type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]Gauge),
		counters: make(map[string]Counter),
	}
}

type GaugeRepositories interface {
	Gauge(name string) (Gauge, bool)
	Gauges() map[string]Gauge
	SetGauge(name string, value float64) error
}

var _ GaugeRepositories = (*MemStorage)(nil)

type CounterRepositories interface {
	Counter(name string) (Counter, bool)
	Counters() map[string]Counter
	AddCounter(name string, value int64) error
}

var _ CounterRepositories = (*MemStorage)(nil)

// Repositories is an interface for managing a set of metrics
type Repositories interface {
	GaugeRepositories
	CounterRepositories
}

var _ Repositories = (*MemStorage)(nil)

// Gauge returns the gauge metric by name
func (m *MemStorage) Gauge(name string) (Gauge, bool) {
	v, ok := m.gauges[name]
	return v, ok
}

// Gauges returns all gauge metrics
func (m *MemStorage) Gauges() map[string]Gauge {
	return m.gauges
}

func (m *MemStorage) SetGauge(name string, value float64) error {
	if m.gauges == nil {
		m.gauges = make(map[string]Gauge)
	}
	g, ok := m.gauges[name]
	if !ok {
		g = Gauge(0)
	}
	err := g.Set(value)
	if err != nil {
		return err
	}
	m.gauges[name] = g
	return nil
}

// Counter returns the counter metric by name
func (m *MemStorage) Counter(name string) (Counter, bool) {
	v, ok := m.counters[name]
	return v, ok
}

// Counters returns all counter metrics
func (m *MemStorage) Counters() map[string]Counter {
	return m.counters
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	if m.counters == nil {
		m.counters = make(map[string]Counter)
	}
	c, ok := m.counters[name]
	if !ok {
		c = Counter(0)
	}
	err := c.Add(value)
	if err != nil {
		return err
	}
	m.counters[name] = c
	return nil
}
