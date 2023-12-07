package collector

import (
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetMetricGauge(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   float64
		data    func() *metrics.MemStorage
		want    func() *metrics.MemStorage
		wantErr bool
	}{
		{
			name:  "Positive case #1",
			key:   "a",
			value: 5.0,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 1.0)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 5.0)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Positive case #2",
			key:   "a",
			value: -5.0,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 1.0)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", -5.0)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Positive case #3",
			key:   "a",
			value: 1.0,
			data: func() *metrics.MemStorage {
				return metrics.NewMemStorage()
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 1.0)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Positive case #4",
			key:   "b",
			value: 3.0,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 1.0)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.SetGauge("a", 1.0)
				_ = data.SetGauge("b", 3.0)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Negative case #1",
			key:   "a",
			value: 5.0,
			data: func() *metrics.MemStorage {
				return nil
			},
			want: func() *metrics.MemStorage {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.data()
			err := setMetricGauge(data, tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want(), data)
		})
	}
}

func TestAddMetricCounter(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   int64
		data    func() *metrics.MemStorage
		want    func() *metrics.MemStorage
		wantErr bool
	}{
		{
			name:  "Positive case #1",
			key:   "a",
			value: 1,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 2)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 3)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Positive case #2",
			key:   "a",
			value: 1,
			data: func() *metrics.MemStorage {
				return metrics.NewMemStorage()
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 1)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Positive case #3",
			key:   "b",
			value: 1,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 2)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 2)
				_ = data.AddCounter("b", 1)
				return data
			},
			wantErr: false,
		},
		{
			name:  "Negative case #1",
			key:   "a",
			value: -1,
			data: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 2)
				return data
			},
			want: func() *metrics.MemStorage {
				data := metrics.NewMemStorage()
				_ = data.AddCounter("a", 2)
				return data
			},
			wantErr: true,
		},
		{
			name:  "Negative case #2",
			key:   "a",
			value: 1,
			data: func() *metrics.MemStorage {
				return nil
			},
			want: func() *metrics.MemStorage {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.data()
			err := addMetricCounter(data, tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want(), data)
		})
	}
}
