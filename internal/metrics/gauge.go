package metrics

type Gauge float64

func (g *Gauge) Set(v float64) error {
	*g = Gauge(v)
	return nil
}
