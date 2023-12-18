package metrics

import "errors"

type Counter int64

func (c *Counter) Add(v int64) error {
	if v < 0 {
		return errors.New(`metrics: the counter value must be positive`)
	}
	*c += Counter(v)
	return nil
}
