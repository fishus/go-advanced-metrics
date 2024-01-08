package metrics

import (
	"io"
	"os"
	"strings"
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

func TestMemStorage_Save(t *testing.T) {
	testCases := []struct {
		name     string
		filename func() string
		storage  *MemStorage
		want     string
		wantErr  bool
	}{
		{
			name: "Positive case #1",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test81278123.json"
				}
				_ = file.Close()
				return file.Name()
			},
			storage: &MemStorage{
				gauges:   map[string]Gauge{"a": {"a", 1.5}},
				counters: map[string]Counter{"b": {"b", 2}},
			},
			want:    `{"gauges":{"a":{"name":"a","value":1.5}},"counters":{"b":{"name":"b","value":2}}}`,
			wantErr: false,
		},
		{
			name: "Positive case #2",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test13640364134.json"
				}
				_ = file.Close()
				return file.Name()
			},
			storage: &MemStorage{
				gauges:   make(map[string]Gauge),
				counters: make(map[string]Counter),
			},
			want:    `{"gauges":{},"counters":{}}`,
			wantErr: false,
		},
		{
			name: "Positive case #3",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test13640364134.json"
				}
				_ = file.Close()
				return file.Name()
			},
			storage: &MemStorage{},
			want:    `{"gauges":null,"counters":null}`,
			wantErr: false,
		},
		{
			name: "Positive case #4",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test902904783.json"
				}
				_ = file.Close()
				os.Remove(file.Name())
				return file.Name()
			},
			storage: &MemStorage{},
			want:    `{"gauges":null,"counters":null}`,
			wantErr: false,
		},
		{
			name: "Negative case #1",
			filename: func() string {
				return ""
			},
			want:    "",
			storage: &MemStorage{},
			wantErr: true,
		},
		{
			name: "Negative case #2",
			filename: func() string {
				return "."
			},
			storage: &MemStorage{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.filename()
			defer os.Remove(filename)

			tc.storage.Filename = filename
			err := tc.storage.Save()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			data, _ := os.ReadFile(filename)

			assert.Equal(t, tc.want, strings.TrimRight(string(data), "\n"))
		})
	}
}

func TestMemStorage_Load(t *testing.T) {
	testCases := []struct {
		name     string
		filename func() string
		data     string
		want     *MemStorage
		wantErr  bool
	}{
		{
			name: "Positive case #1",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test712390123.json"
				}
				_ = file.Close()

				return file.Name()
			},
			data: `{"gauges":{"a":{"name":"a","value":1.5}},"counters":{"b":{"name":"b","value":2}}}`,
			want: &MemStorage{
				gauges:   map[string]Gauge{"a": {"a", 1.5}},
				counters: map[string]Counter{"b": {"b", 2}},
			},
			wantErr: false,
		},
		{
			name: "Positive case #2",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test6123612.json"
				}
				_ = file.Close()

				return file.Name()
			},
			data: `{"gauges":null,"counters":null}`,
			want: &MemStorage{
				gauges:   map[string]Gauge(nil),
				counters: map[string]Counter(nil),
			},
			wantErr: false,
		},
		{
			name: "Positive case #3",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test6123612.json"
				}
				_ = file.Close()

				return file.Name()
			},
			data: `{"gauges":{},"counters":{}}`,
			want: &MemStorage{
				gauges:   map[string]Gauge{},
				counters: map[string]Counter{},
			},
			wantErr: false,
		},
		{
			name: "Positive case #4",
			filename: func() string {
				tmpDir := os.TempDir()
				file, err := os.CreateTemp(tmpDir, "test*.json")
				if err != nil {
					return "test6109237.json"
				}
				_ = file.Close()

				return file.Name()
			},
			data: "",
			want: &MemStorage{
				gauges:   map[string]Gauge{},
				counters: map[string]Counter{},
			},
			wantErr: false,
		},
		{
			name: "Negative case #1",
			filename: func() string {
				return ""
			},
			want: &MemStorage{
				gauges:   map[string]Gauge{},
				counters: map[string]Counter{},
			},
			wantErr: true,
		},
		{
			name: "Negative case #2",
			filename: func() string {
				return "."
			},
			want: &MemStorage{
				gauges:   map[string]Gauge{},
				counters: map[string]Counter{},
			},
			wantErr: true,
		},
		{
			name: "Negative case #3",
			filename: func() string {
				return "."
			},
			data: `{"gauges":{}`,
			want: &MemStorage{
				gauges:   map[string]Gauge{},
				counters: map[string]Counter{},
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.filename()
			tc.want.Filename = filename
			defer os.Remove(filename)

			if filename != "" && tc.data != "" {
				_ = os.WriteFile(filename, []byte(tc.data), 0664)
			}

			storage := NewMemStorage()
			storage.Filename = filename

			err := storage.Load()
			if tc.wantErr {
				require.Error(t, err)
			} else if err != nil {
				assert.ErrorIs(t, err, io.EOF)
			} else {
				require.NoError(t, err)
			}

			assert.EqualValues(t, tc.want, storage)
		})
	}
}
