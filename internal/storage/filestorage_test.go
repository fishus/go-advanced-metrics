package storage

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

func TestNewFileStorage(t *testing.T) {
	want := &FileStorage{filename: "file.txt"}
	want.gauges = make(map[string]metrics.Gauge)
	want.counters = make(map[string]metrics.Counter)
	got := NewFileStorage("file.txt")
	assert.Equal(t, want, got)
}

func TestFileStorage_Save(t *testing.T) {
	testCases := []struct {
		name     string
		filename func() string
		storage  *FileStorage
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
			storage: func() *FileStorage {
				gauges := map[string]metrics.Gauge{}
				ga, _ := metrics.NewGauge("a", 1.5)
				gauges["a"] = *ga

				counters := map[string]metrics.Counter{}
				cb, _ := metrics.NewCounter("b", 2)
				counters["b"] = *cb

				fs := &FileStorage{}
				fs.gauges = gauges
				fs.counters = counters
				return fs
			}(),
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
			storage: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = make(map[string]metrics.Gauge)
				fs.counters = make(map[string]metrics.Counter)
				return fs
			}(),
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
			storage: &FileStorage{},
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
			storage: &FileStorage{},
			want:    `{"gauges":null,"counters":null}`,
			wantErr: false,
		},
		{
			name: "Negative case #1",
			filename: func() string {
				return ""
			},
			want:    "",
			storage: &FileStorage{},
			wantErr: true,
		},
		{
			name: "Negative case #2",
			filename: func() string {
				return "."
			},
			storage: &FileStorage{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.filename()
			defer os.Remove(filename)

			tc.storage.filename = filename
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

func TestFileStorage_Load(t *testing.T) {
	testCases := []struct {
		name     string
		filename func() string
		data     string
		want     *FileStorage
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
			want: func() *FileStorage {
				gauges := map[string]metrics.Gauge{}
				ga, _ := metrics.NewGauge("a", 1.5)
				gauges["a"] = *ga

				counters := map[string]metrics.Counter{}
				cb, _ := metrics.NewCounter("b", 2)
				counters["b"] = *cb

				fs := &FileStorage{}
				fs.gauges = gauges
				fs.counters = counters
				return fs
			}(),
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
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge(nil)
				fs.counters = map[string]metrics.Counter(nil)
				return fs
			}(),
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
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge{}
				fs.counters = map[string]metrics.Counter{}
				return fs
			}(),
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
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge{}
				fs.counters = map[string]metrics.Counter{}
				return fs
			}(),
			wantErr: false,
		},
		{
			name: "Negative case #1",
			filename: func() string {
				return ""
			},
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge{}
				fs.counters = map[string]metrics.Counter{}
				return fs
			}(),
			wantErr: true,
		},
		{
			name: "Negative case #2",
			filename: func() string {
				return "."
			},
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge{}
				fs.counters = map[string]metrics.Counter{}
				return fs
			}(),
			wantErr: true,
		},
		{
			name: "Negative case #3",
			filename: func() string {
				return "."
			},
			data: `{"gauges":{}`,
			want: func() *FileStorage {
				fs := &FileStorage{}
				fs.gauges = map[string]metrics.Gauge{}
				fs.counters = map[string]metrics.Counter{}
				return fs
			}(),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.filename()
			tc.want.filename = filename
			defer os.Remove(filename)

			if filename != "" && tc.data != "" {
				_ = os.WriteFile(filename, []byte(tc.data), 0664)
			}

			storage := NewFileStorage(filename)

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
