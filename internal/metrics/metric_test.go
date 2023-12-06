package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMetric(t *testing.T) {
	want := Metric{}
	got := NewMetric()
	assert.Equal(t, want, got)
}

func TestMetric_SetGauge(t *testing.T) {
	testCases := []struct {
		name    string
		value   float64
		metric  Metric
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			value:   2.2,
			metric:  Metric{gauge: 1.0},
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			value:   -1.5,
			metric:  Metric{gauge: 1.0},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			value:   0,
			metric:  Metric{gauge: 1.0},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.metric.SetGauge(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.value, tc.metric.gauge)
		})
	}
}

func TestMetric_AddCounter(t *testing.T) {
	testCases := []struct {
		name    string
		value   int64
		metric  Metric
		want    int64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			value:   1,
			metric:  Metric{counter: 2},
			want:    3,
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			value:   0,
			metric:  Metric{counter: 1},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			value:   -1,
			metric:  Metric{counter: 2},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.metric.AddCounter(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, tc.metric.counter)
		})
	}
}

func TestMetric_Gauge(t *testing.T) {
	testCases := []struct {
		name   string
		metric Metric
	}{
		{
			name:   "Positive case #1",
			metric: Metric{gauge: 1.23},
		},
		{
			name:   "Positive case #2",
			metric: Metric{gauge: -1.23},
		},
		{
			name:   "Positive case #3",
			metric: Metric{gauge: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.metric.gauge, tc.metric.Gauge())
		})
	}
}

func TestMetric_Counter(t *testing.T) {
	testCases := []struct {
		name   string
		metric Metric
	}{
		{
			name:   "Positive case #1",
			metric: Metric{counter: 1},
		},
		{
			name:   "Positive case #1",
			metric: Metric{counter: 0},
		},
		{
			name:   "Positive case #1",
			metric: Metric{counter: 100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.metric.counter, tc.metric.Counter())
		})
	}
}
