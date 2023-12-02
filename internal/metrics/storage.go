package metrics

// MemStorage contains a set of values for all metrics
type MemStorage struct {
	metrics map[string]Metric
}

// Storager is an interface for managing a set of metrics
type Storager interface {
	SetGauge(name string, value float64) error
	AddCounter(name string, value int64) error
	Metric(name string) (Metric, bool)
	Metrics() map[string]Metric
}

var _ Storager = (*MemStorage)(nil)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]Metric),
	}
}

func (m *MemStorage) SetGauge(name string, value float64) error {
	if m.metrics == nil {
		m.metrics = map[string]Metric{}
	}
	metric, ok := m.metrics[name]
	if !ok {
		metric = NewMetric()
	}
	err := metric.SetGauge(value)
	if err != nil {
		return err
	}
	m.metrics[name] = metric
	return nil
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	if m.metrics == nil {
		m.metrics = map[string]Metric{}
	}
	metric, ok := m.metrics[name]
	if !ok {
		metric = NewMetric()
	}
	err := metric.AddCounter(value)
	if err != nil {
		return err
	}
	m.metrics[name] = metric
	return nil
}

func (m *MemStorage) Metric(name string) (Metric, bool) {
	v, ok := m.metrics[name]
	return v, ok
}

func (m *MemStorage) Metrics() map[string]Metric {
	return m.metrics
}
