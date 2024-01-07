package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemStorage(t *testing.T) {
	want := &MemStorage{
		gauges:   map[string]Gauge{},
		counters: map[string]Counter{},
	}
	got := NewMemStorage()
	assert.Equal(t, want, got)
}

func TestMemStorage_Gauge(t *testing.T) {
	type want struct {
		gauge Gauge
		ok    bool
	}
	testCases := []struct {
		name   string
		gauges map[string]Gauge
		key    string
		want   want
	}{
		{
			name:   "Positive case #1",
			gauges: map[string]Gauge{"a": {"a", 2.1}, "b": {"b", -1.5}},
			key:    "b",
			want: want{
				gauge: Gauge{"b", -1.5},
				ok:    true,
			},
		},
		{
			name:   "Negative case #1",
			gauges: map[string]Gauge{},
			key:    "a",
			want: want{
				ok: false,
			},
		},
		{
			name:   "Negative case #2",
			gauges: map[string]Gauge{"a": {"a", 2.1}},
			key:    "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				gauges: tc.gauges,
			}
			g, ok := m.Gauge(tc.key)
			require.Equal(t, tc.want.ok, ok)
			if tc.want.ok {
				assert.EqualValues(t, tc.want.gauge, g)
			}
		})
	}
}

func TestMemStorage_GaugeValue(t *testing.T) {
	type want struct {
		value float64
		ok    bool
	}
	testCases := []struct {
		name   string
		gauges map[string]Gauge
		key    string
		want   want
	}{
		{
			name:   "Positive case #1",
			gauges: map[string]Gauge{"a": {"a", 2.1}, "b": {"b", -1.5}},
			key:    "b",
			want: want{
				value: -1.5,
				ok:    true,
			},
		},
		{
			name:   "Negative case #1",
			gauges: map[string]Gauge{},
			key:    "a",
			want: want{
				ok: false,
			},
		},
		{
			name:   "Negative case #2",
			gauges: map[string]Gauge{"a": {"a", 2.1}},
			key:    "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				gauges: tc.gauges,
			}
			g, ok := m.GaugeValue(tc.key)
			if !tc.want.ok {
				assert.Equal(t, tc.want.ok, ok)
				return
			}
			require.Equal(t, tc.want.ok, ok)
			assert.Equal(t, tc.want.value, g)
		})
	}
}

func TestMemStorage_Gauges(t *testing.T) {
	gauges := map[string]Gauge{}
	gauges["a"] = Gauge{"a", 1.0}
	gauges["b"] = Gauge{"b", 2.1}

	m := &MemStorage{
		gauges: gauges,
	}
	assert.Equal(t, gauges, m.Gauges())
}

func TestMemStorage_SetGauge(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   float64
		gauges  map[string]Gauge
		want    map[string]Gauge
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			key:     "a",
			value:   5.0,
			gauges:  map[string]Gauge{"a": {"a", 1.0}},
			want:    map[string]Gauge{"a": {"a", 5.0}},
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			key:     "a",
			value:   -5.0,
			gauges:  map[string]Gauge{"a": {"a", 1.0}},
			want:    map[string]Gauge{"a": {"a", -5.0}},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "a",
			value:   1.0,
			gauges:  map[string]Gauge{},
			want:    map[string]Gauge{"a": {"a", 1.0}},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "b",
			value:   3.0,
			gauges:  map[string]Gauge{"a": {"a", 1.0}},
			want:    map[string]Gauge{"a": {"a", 1.0}, "b": {"b", 3.0}},
			wantErr: false,
		},
		{
			name:    "Positive case #5",
			key:     "a",
			value:   5.0,
			gauges:  nil,
			want:    map[string]Gauge{"a": {"a", 5}},
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			key:     "",
			value:   5.0,
			gauges:  map[string]Gauge{},
			want:    map[string]Gauge{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				gauges: tc.gauges,
			}
			err := m.SetGauge(tc.key, tc.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, m.gauges)
		})
	}
}

func TestMemStorage_Counter(t *testing.T) {
	type want struct {
		counter Counter
		ok      bool
	}
	testCases := []struct {
		name     string
		counters map[string]Counter
		key      string
		want     want
	}{
		{
			name:     "Positive case #1",
			counters: map[string]Counter{"a": {"a", 10}, "b": {"b", 20}},
			key:      "b",
			want: want{
				counter: Counter{"b", 20},
				ok:      true,
			},
		},
		{
			name:     "Negative case #1",
			counters: map[string]Counter{},
			key:      "a",
			want: want{
				ok: false,
			},
		},
		{
			name:     "Negative case #2",
			counters: map[string]Counter{"a": {"a", 10}},
			key:      "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				counters: tc.counters,
			}
			c, ok := m.Counter(tc.key)
			require.Equal(t, tc.want.ok, ok)
			if tc.want.ok {
				assert.EqualValues(t, tc.want.counter, c)
			}
		})
	}
}

func TestMemStorage_CounterValue(t *testing.T) {
	type want struct {
		value int64
		ok    bool
	}
	testCases := []struct {
		name     string
		counters map[string]Counter
		key      string
		want     want
	}{
		{
			name:     "Positive case #1",
			counters: map[string]Counter{"a": {"a", 10}, "b": {"b", 20}},
			key:      "b",
			want: want{
				value: 20,
				ok:    true,
			},
		},
		{
			name:     "Negative case #1",
			counters: map[string]Counter{},
			key:      "a",
			want: want{
				ok: false,
			},
		},
		{
			name:     "Negative case #2",
			counters: map[string]Counter{"a": {"a", 10}},
			key:      "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				counters: tc.counters,
			}
			c, ok := m.CounterValue(tc.key)
			if !tc.want.ok {
				assert.Equal(t, tc.want.ok, ok)
				return
			}
			require.Equal(t, tc.want.ok, ok)
			assert.Equal(t, tc.want.value, c)
		})
	}
}

func TestMemStorage_Counters(t *testing.T) {
	counters := map[string]Counter{}
	counters["a"] = Counter{"a", 1}
	counters["b"] = Counter{"b", 100}

	m := &MemStorage{
		counters: counters,
	}
	assert.Equal(t, counters, m.Counters())
}

func TestMemStorage_AddCounter(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		value    int64
		counters map[string]Counter
		want     map[string]Counter
		wantErr  bool
	}{
		{
			name:     "Positive case #1",
			key:      "a",
			value:    1,
			counters: map[string]Counter{"a": {"a", 2}},
			want:     map[string]Counter{"a": {"a", 3}},
			wantErr:  false,
		},
		{
			name:     "Positive case #2",
			key:      "a",
			value:    1,
			counters: map[string]Counter{},
			want:     map[string]Counter{"a": {"a", 1}},
			wantErr:  false,
		},
		{
			name:     "Positive case #3",
			key:      "b",
			value:    1,
			counters: map[string]Counter{"a": {"a", 2}},
			want:     map[string]Counter{"a": {"a", 2}, "b": {"b", 1}},
			wantErr:  false,
		},
		{
			name:     "Positive case #4",
			key:      "a",
			value:    1,
			counters: nil,
			want:     map[string]Counter{"a": {"a", 1}},
			wantErr:  false,
		},
		{
			name:     "Negative case #1",
			key:      "a",
			value:    -1,
			counters: map[string]Counter{"a": {"a", 2}},
			want:     map[string]Counter{"a": {"a", 2}},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &MemStorage{
				counters: tc.counters,
			}
			err := m.AddCounter(tc.key, tc.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, m.counters)
		})
	}
}
