package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	want := &MemStorage{
		metrics: metrics{},
	}
	got := NewMemStorage()
	assert.Equal(t, want, got)
}

func TestMemStorage_Metrics(t *testing.T) {
	metrics := metrics{}
	metrics["test"] = Metric{gauge: 1.0, counter: 10}

	ms := &MemStorage{
		metrics: metrics,
	}
	assert.Equal(t, metrics, ms.Metrics())
}

func TestMemStorage_Metric(t *testing.T) {
	type want struct {
		metric Metric
		ok     bool
	}
	testCases := []struct {
		name    string
		metrics metrics
		key     string
		want    want
	}{
		{
			name:    "Positive case #1",
			metrics: metrics{"test": Metric{gauge: 1.0, counter: 10}},
			key:     "test",
			want: want{
				metric: Metric{gauge: 1.0, counter: 10},
				ok:     true,
			},
		},
		{
			name:    "Negative case #1",
			metrics: metrics{},
			key:     "test",
			want: want{
				ok: false,
			},
		},
		{
			name:    "Negative case #2",
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 10}},
			key:     "bbb",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := &MemStorage{
				metrics: tc.metrics,
			}
			m, ok := ms.Metric(tc.key)
			if !tc.want.ok {
				assert.Equal(t, tc.want.ok, ok)
				return
			}
			require.Equal(t, tc.want.ok, ok)
			assert.Equal(t, tc.want.metric, m)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   float64
		metrics metrics
		want    metrics
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			key:     "aaa",
			value:   5.0,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 10}},
			want:    metrics{"aaa": Metric{gauge: 5.0, counter: 10}},
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			key:     "aaa",
			value:   -5.0,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 10}},
			want:    metrics{"aaa": Metric{gauge: -5.0, counter: 10}},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "aaa",
			value:   1.0,
			metrics: metrics{},
			want:    metrics{"aaa": Metric{gauge: 1.0, counter: 0}},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "bbb",
			value:   3.0,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 2}},
			want:    metrics{"aaa": Metric{gauge: 1.0, counter: 2}, "bbb": Metric{gauge: 3.0, counter: 0}},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := &MemStorage{
				metrics: tc.metrics,
			}
			err := ms.SetGauge(tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, ms.metrics)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   int64
		metrics metrics
		want    metrics
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			key:     "aaa",
			value:   1,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 2}},
			want:    metrics{"aaa": Metric{gauge: 1.0, counter: 3}},
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			key:     "aaa",
			value:   1,
			metrics: metrics{},
			want:    metrics{"aaa": Metric{gauge: 0, counter: 1}},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "bbb",
			value:   1,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 2}},
			want:    metrics{"aaa": Metric{gauge: 1.0, counter: 2}, "bbb": Metric{gauge: 0, counter: 1}},
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			key:     "aaa",
			value:   -1,
			metrics: metrics{"aaa": Metric{gauge: 1.0, counter: 2}},
			want:    metrics{"aaa": Metric{gauge: 1.0, counter: 2}},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := &MemStorage{
				metrics: tc.metrics,
			}
			err := ms.AddCounter(tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, ms.metrics)
		})
	}
}
