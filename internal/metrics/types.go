package metrics

type MetricType string

const (
	TypeGauge   MetricType = "gauge"
	TypeCounter MetricType = "counter"
)
