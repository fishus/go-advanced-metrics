package main

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
		metrics func() map[string]metrics.Metric
		want    func() map[string]metrics.Metric
		wantErr bool
	}{
		{
			name:  "Positive case #1",
			key:   "a",
			value: 5.0,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(1)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(5)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: false,
		},
		{
			name:  "Positive case #2",
			key:   "a",
			value: -5.0,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(1)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(-5)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: false,
		},
		{
			name:  "Positive case #3",
			key:   "a",
			value: 1.0,
			metrics: func() map[string]metrics.Metric {
				return map[string]metrics.Metric{}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(1)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: false,
		},
		{
			name:  "Positive case #4",
			key:   "b",
			value: 3.0,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(1)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.SetGauge(1)
				b := metrics.NewMetric()
				_ = b.SetGauge(3)
				return map[string]metrics.Metric{"a": a, "b": b}
			},
			wantErr: false,
		},
		{
			name:  "Negative case #1",
			key:   "a",
			value: 5.0,
			metrics: func() map[string]metrics.Metric {
				return nil
			},
			want: func() map[string]metrics.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mtx := tc.metrics()
			err := setMetricGauge(mtx, tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want(), mtx)
		})
	}
}

func TestAddMetricCounter(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   int64
		metrics func() map[string]metrics.Metric
		want    func() map[string]metrics.Metric
		wantErr bool
	}{
		{
			name:  "Positive case #1",
			key:   "a",
			value: 1,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(2)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(3)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: false,
		},
		{
			name:  "Positive case #2",
			key:   "a",
			value: 1,
			metrics: func() map[string]metrics.Metric {
				return map[string]metrics.Metric{}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(1)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: false,
		},
		{
			name:  "Positive case #3",
			key:   "b",
			value: 1,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(2)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(2)
				b := metrics.NewMetric()
				_ = b.AddCounter(1)
				return map[string]metrics.Metric{"a": a, "b": b}
			},
			wantErr: false,
		},
		{
			name:  "Negative case #1",
			key:   "a",
			value: -1,
			metrics: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(2)
				return map[string]metrics.Metric{"a": a}
			},
			want: func() map[string]metrics.Metric {
				a := metrics.NewMetric()
				_ = a.AddCounter(2)
				return map[string]metrics.Metric{"a": a}
			},
			wantErr: true,
		},
		{
			name:  "Negative case #2",
			key:   "a",
			value: 1,
			metrics: func() map[string]metrics.Metric {
				return nil
			},
			want: func() map[string]metrics.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mtx := tc.metrics()
			err := addMetricCounter(mtx, tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want(), mtx)
		})
	}
}
