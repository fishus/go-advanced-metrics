package metrics

import (
	"encoding/json"
	"errors"
)

// Counter implements the metric type Counter.
type Counter struct {
	name  string
	value int64
}

// NewCounter returns a pointer to the Counter structure.
func NewCounter(name string, v int64) (*Counter, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	counter := &Counter{name: name}
	if err := counter.AddValue(v); err != nil {
		return nil, err
	}
	return counter, nil
}

// Name returns the name of the counter.
func (c Counter) Name() string {
	return c.name
}

// Value returns the counter value.
func (c Counter) Value() int64 {
	return c.value
}

// AddValue increases the counter value.
func (c *Counter) AddValue(v int64) error {
	if v < 0 {
		return errors.New(`metrics: the counter value must be positive`)
	}
	c.value += v
	return nil
}

// MarshalJSON implements the Marshaler interface.
func (c Counter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}{
		Name:  c.name,
		Value: c.value,
	})
}

// UnmarshalJSON implements the Unmarshaler interface.
func (c *Counter) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	c.name = aux.Name
	c.value = aux.Value

	return nil
}
