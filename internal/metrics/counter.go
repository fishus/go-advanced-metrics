package metrics

import "errors"

type Counter struct {
	Name  string
	Value int64
}

func NewCounter(name string, v int64) (*Counter, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	counter := &Counter{Name: name}
	if err := counter.AddValue(v); err != nil {
		return nil, err
	}
	return counter, nil
}

func (c *Counter) AddValue(v int64) error {
	if v < 0 {
		return errors.New(`metrics: the counter value must be positive`)
	}
	c.Value += v
	return nil
}
