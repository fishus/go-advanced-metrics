package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

func TestNewMemStorage(t *testing.T) {
	want := &MemStorage{
		gauges:   map[string]metrics.Gauge{},
		counters: map[string]metrics.Counter{},
	}
	got := NewMemStorage()
	assert.Equal(t, want, got)
}

func TestMemStorage_Gauge(t *testing.T) {
	type gauge struct {
		name  string
		value float64
	}

	type want struct {
		gauge metrics.Gauge
		ok    bool
	}

	testCases := []struct {
		name   string
		gauges []gauge
		key    string
		want   want
	}{
		{
			name: "Positive case #1",
			gauges: []gauge{
				{"a", 2.1},
				{"b", -1.5},
			},
			key: "b",
			want: want{
				gauge: func() metrics.Gauge {
					g, _ := metrics.NewGauge("b", -1.5)
					return *g
				}(),
				ok: true,
			},
		},
		{
			name:   "Negative case #1",
			gauges: []gauge{},
			key:    "a",
			want: want{
				ok: false,
			},
		},
		{
			name:   "Negative case #2",
			gauges: []gauge{{"a", 2.1}},
			key:    "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauges := map[string]metrics.Gauge{}
			if tc.gauges != nil {
				for _, v := range tc.gauges {
					g, _ := metrics.NewGauge(v.name, v.value)
					gauges[v.name] = *g
				}
			} else {
				gauges = nil
			}
			m := &MemStorage{
				gauges: gauges,
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
	type gauge struct {
		name  string
		value float64
	}

	type want struct {
		value float64
		ok    bool
	}

	testCases := []struct {
		name   string
		gauges []gauge
		key    string
		want   want
	}{
		{
			name: "Positive case #1",
			gauges: []gauge{
				{"a", 2.1},
				{"b", -1.5},
			},
			key: "b",
			want: want{
				value: -1.5,
				ok:    true,
			},
		},
		{
			name:   "Negative case #1",
			gauges: []gauge{},
			key:    "a",
			want: want{
				ok: false,
			},
		},
		{
			name:   "Negative case #2",
			gauges: []gauge{{"a", 2.1}},
			key:    "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauges := map[string]metrics.Gauge{}
			if tc.gauges != nil {
				for _, v := range tc.gauges {
					g, _ := metrics.NewGauge(v.name, v.value)
					gauges[v.name] = *g
				}
			} else {
				gauges = nil
			}
			m := &MemStorage{
				gauges: gauges,
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
	gauges := map[string]metrics.Gauge{}
	a, _ := metrics.NewGauge("a", 1.0)
	b, _ := metrics.NewGauge("b", 2.1)
	gauges["a"] = *a
	gauges["b"] = *b

	m := &MemStorage{
		gauges: gauges,
	}
	assert.Equal(t, gauges, m.Gauges())
}

func TestMemStorage_GaugesFiltered(t *testing.T) {
	gauges := map[string]metrics.Gauge{}
	a, _ := metrics.NewGauge("a", 1.0)
	b, _ := metrics.NewGauge("b", 2.1)
	c, _ := metrics.NewGauge("c", 3.4)
	gauges["a"] = *a
	gauges["b"] = *b
	gauges["c"] = *c

	m := &MemStorage{
		gauges: gauges,
	}

	filter := []string{"b", "c"}

	want := map[string]metrics.Gauge{}
	want["b"] = *b
	want["c"] = *c

	assert.Equal(t, want, m.Gauges(FilterNames(filter)))
}

func TestMemStorage_SetGauge(t *testing.T) {
	type gauge struct {
		name  string
		value float64
	}

	testCases := []struct {
		name    string
		key     string
		value   float64
		gauges  []gauge
		want    map[string]metrics.Gauge
		wantErr bool
	}{
		{
			name:   "Positive case #1",
			key:    "a",
			value:  5.0,
			gauges: []gauge{{"a", 1.0}},
			want: func() map[string]metrics.Gauge {
				g := map[string]metrics.Gauge{}
				a, _ := metrics.NewGauge("a", 5.0)
				g["a"] = *a
				return g
			}(),
			wantErr: false,
		},
		{
			name:   "Positive case #2",
			key:    "a",
			value:  -5.0,
			gauges: []gauge{{"a", 1.0}},
			want: func() map[string]metrics.Gauge {
				g := map[string]metrics.Gauge{}
				a, _ := metrics.NewGauge("a", -5.0)
				g["a"] = *a
				return g
			}(),
			wantErr: false,
		},
		{
			name:   "Positive case #3",
			key:    "a",
			value:  1.0,
			gauges: []gauge{},
			want: func() map[string]metrics.Gauge {
				g := map[string]metrics.Gauge{}
				a, _ := metrics.NewGauge("a", 1.0)
				g["a"] = *a
				return g
			}(),
			wantErr: false,
		},
		{
			name:   "Positive case #4",
			key:    "b",
			value:  3.0,
			gauges: []gauge{{"a", 1.0}},
			want: func() map[string]metrics.Gauge {
				g := map[string]metrics.Gauge{}
				a, _ := metrics.NewGauge("a", 1.0)
				g["a"] = *a
				b, _ := metrics.NewGauge("b", 3.0)
				g["b"] = *b
				return g
			}(),
			wantErr: false,
		},
		{
			name:   "Positive case #5",
			key:    "a",
			value:  5.0,
			gauges: nil,
			want: func() map[string]metrics.Gauge {
				g := map[string]metrics.Gauge{}
				a, _ := metrics.NewGauge("a", 5)
				g["a"] = *a
				return g
			}(),
			wantErr: false,
		},
		{
			name:   "Negative case #1",
			key:    "",
			value:  5.0,
			gauges: []gauge{},
			want: func() map[string]metrics.Gauge {
				return map[string]metrics.Gauge{}
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauges := map[string]metrics.Gauge{}
			if tc.gauges != nil {
				for _, v := range tc.gauges {
					g, _ := metrics.NewGauge(v.name, v.value)
					gauges[v.name] = *g
				}
			} else {
				gauges = nil
			}
			m := &MemStorage{
				gauges: gauges,
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

func BenchmarkMemStorage_SetGauge(b *testing.B) {
	m := &MemStorage{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.SetGauge("a", 123.45)
	}
}

func TestMemStorage_Counter(t *testing.T) {
	type counter struct {
		name  string
		value int64
	}

	type want struct {
		counter metrics.Counter
		ok      bool
	}

	testCases := []struct {
		name     string
		counters []counter
		key      string
		want     want
	}{
		{
			name:     "Positive case #1",
			counters: []counter{{"a", 10}, {"b", 20}},
			key:      "b",
			want: want{
				counter: func() metrics.Counter {
					c, _ := metrics.NewCounter("b", 20)
					return *c
				}(),
				ok: true,
			},
		},
		{
			name:     "Negative case #1",
			counters: []counter{},
			key:      "a",
			want: want{
				ok: false,
			},
		},
		{
			name:     "Negative case #2",
			counters: []counter{{"a", 10}},
			key:      "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[string]metrics.Counter{}
			if tc.counters != nil {
				for _, v := range tc.counters {
					c, _ := metrics.NewCounter(v.name, v.value)
					counters[v.name] = *c
				}
			} else {
				counters = nil
			}
			m := &MemStorage{
				counters: counters,
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
	type counter struct {
		name  string
		value int64
	}

	type want struct {
		value int64
		ok    bool
	}

	testCases := []struct {
		name     string
		counters []counter
		key      string
		want     want
	}{
		{
			name:     "Positive case #1",
			counters: []counter{{"a", 10}, {"b", 20}},
			key:      "b",
			want: want{
				value: 20,
				ok:    true,
			},
		},
		{
			name:     "Negative case #1",
			counters: []counter{},
			key:      "a",
			want: want{
				ok: false,
			},
		},
		{
			name:     "Negative case #2",
			counters: []counter{{"a", 10}},
			key:      "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[string]metrics.Counter{}
			if tc.counters != nil {
				for _, v := range tc.counters {
					c, _ := metrics.NewCounter(v.name, v.value)
					counters[v.name] = *c
				}
			} else {
				counters = nil
			}
			m := &MemStorage{
				counters: counters,
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
	counters := map[string]metrics.Counter{}
	a, _ := metrics.NewCounter("a", 1)
	b, _ := metrics.NewCounter("b", 100)
	counters["a"] = *a
	counters["b"] = *b

	m := &MemStorage{
		counters: counters,
	}
	assert.Equal(t, counters, m.Counters())
}

func TestMemStorage_CountersFiltered(t *testing.T) {
	counters := map[string]metrics.Counter{}
	a, _ := metrics.NewCounter("a", 1)
	b, _ := metrics.NewCounter("b", 100)
	c, _ := metrics.NewCounter("c", 1000)
	counters["a"] = *a
	counters["b"] = *b
	counters["c"] = *c

	m := &MemStorage{
		counters: counters,
	}

	filter := []string{"b", "c"}

	want := map[string]metrics.Counter{}
	want["b"] = *b
	want["c"] = *c

	assert.Equal(t, want, m.Counters(FilterNames(filter)))
}

func TestMemStorage_AddCounter(t *testing.T) {
	type counter struct {
		name  string
		value int64
	}

	testCases := []struct {
		name     string
		key      string
		value    int64
		counters []counter
		want     map[string]metrics.Counter
		wantErr  bool
	}{
		{
			name:     "Positive case #1",
			key:      "a",
			value:    1,
			counters: []counter{{"a", 2}},
			want: func() map[string]metrics.Counter {
				c := map[string]metrics.Counter{}
				a, _ := metrics.NewCounter("a", 3)
				c["a"] = *a
				return c
			}(),
			wantErr: false,
		},
		{
			name:     "Positive case #2",
			key:      "a",
			value:    1,
			counters: []counter{},
			want: func() map[string]metrics.Counter {
				c := map[string]metrics.Counter{}
				a, _ := metrics.NewCounter("a", 1)
				c["a"] = *a
				return c
			}(),
			wantErr: false,
		},
		{
			name:     "Positive case #3",
			key:      "b",
			value:    1,
			counters: []counter{{"a", 2}},
			want: func() map[string]metrics.Counter {
				c := map[string]metrics.Counter{}
				a, _ := metrics.NewCounter("a", 2)
				c["a"] = *a
				b, _ := metrics.NewCounter("b", 1)
				c["b"] = *b
				return c
			}(),
			wantErr: false,
		},
		{
			name:     "Positive case #4",
			key:      "a",
			value:    1,
			counters: nil,
			want: func() map[string]metrics.Counter {
				c := map[string]metrics.Counter{}
				a, _ := metrics.NewCounter("a", 1)
				c["a"] = *a
				return c
			}(),
			wantErr: false,
		},
		{
			name:     "Negative case #1",
			key:      "a",
			value:    -1,
			counters: []counter{{"a", 2}},
			want: func() map[string]metrics.Counter {
				c := map[string]metrics.Counter{}
				a, _ := metrics.NewCounter("a", 2)
				c["a"] = *a
				return c
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[string]metrics.Counter{}
			if tc.counters != nil {
				for _, v := range tc.counters {
					c, _ := metrics.NewCounter(v.name, v.value)
					counters[v.name] = *c
				}
			} else {
				counters = nil
			}
			m := &MemStorage{
				counters: counters,
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

func TestMemStorage_InsertBatch(t *testing.T) {
	type counter struct {
		name  string
		value int64
	}

	type gauge struct {
		name  string
		value float64
	}

	testCases := []struct {
		name     string
		storage  *MemStorage
		counters []counter
		gauges   []gauge
		want     *MemStorage
	}{
		{
			name: "Insert counters",
			counters: []counter{
				{"a", 2},
				{"b", 3},
			},
			storage: func() *MemStorage {
				m := &MemStorage{}
				m.AddCounter("a", 5)
				return m
			}(),
			want: func() *MemStorage {
				m := &MemStorage{}
				m.AddCounter("a", 7)
				m.AddCounter("b", 3)
				return m
			}(),
		},
		{
			name: "Insert gauges",
			gauges: []gauge{
				{name: "a", value: 1.2},
				{name: "b", value: 2.3},
			},
			storage: func() *MemStorage {
				m := &MemStorage{}
				m.SetGauge("a", 5)
				return m
			}(),
			want: func() *MemStorage {
				m := &MemStorage{}
				m.SetGauge("a", 1.2)
				m.SetGauge("b", 2.3)
				return m
			}(),
		},
		{
			name: "Insert counters and gauges",
			counters: []counter{
				{"a", 2},
				{"b", 3},
			},
			gauges: []gauge{
				{name: "a", value: 1.2},
				{name: "b", value: 2.3},
			},
			storage: func() *MemStorage {
				m := &MemStorage{}
				m.SetGauge("a", 5)
				m.AddCounter("a", 5)
				return m
			}(),
			want: func() *MemStorage {
				m := &MemStorage{}
				m.SetGauge("a", 1.2)
				m.SetGauge("b", 2.3)
				m.AddCounter("a", 7)
				m.AddCounter("b", 3)
				return m
			}(),
		},
		{
			name: "Insert nothing",
			storage: func() *MemStorage {
				m := &MemStorage{}
				return m
			}(),
			want: func() *MemStorage {
				m := &MemStorage{}
				return m
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.counters) == 0 && len(tc.gauges) == 0 {
				err := tc.storage.InsertBatch()
				require.NoError(t, err)
				return
			}

			var countersBatch []metrics.Counter
			if len(tc.counters) > 0 {
				for _, v := range tc.counters {
					c, err := metrics.NewCounter(v.name, v.value)
					if err == nil {
						countersBatch = append(countersBatch, *c)
					}
				}
			}

			var gaugesBatch []metrics.Gauge
			if len(tc.gauges) > 0 {
				for _, v := range tc.gauges {
					g, err := metrics.NewGauge(v.name, v.value)
					if err == nil {
						gaugesBatch = append(gaugesBatch, *g)
					}
				}
			}

			err := tc.storage.InsertBatch(WithCounters(countersBatch), WithGauges(gaugesBatch))
			require.NoError(t, err)
			assert.EqualValues(t, tc.want.Counters(), tc.storage.Counters())
			assert.EqualValues(t, tc.want.Gauges(), tc.storage.Gauges())
		})
	}
}
