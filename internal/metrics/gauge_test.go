package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGauge(t *testing.T) {
	type attr struct {
		name  string
		value float64
	}

	testCases := []struct {
		name    string
		attr    attr
		want    Gauge
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			attr:    attr{name: "test", value: 2.2},
			want:    Gauge{name: "test", value: 2.2},
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			attr:    attr{name: "", value: 1.0},
			want:    Gauge{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauge, err := NewGauge(tc.attr.name, tc.attr.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tc.want, *gauge)
		})
	}
}

func TestGauge_SetValue(t *testing.T) {
	type attr struct {
		name  string
		value float64
	}

	testCases := []struct {
		name    string
		attr    attr
		want    Gauge
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			attr:    attr{name: "test", value: 2.2},
			want:    Gauge{name: "test", value: 2.2},
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			attr:    attr{name: "test", value: -1.5},
			want:    Gauge{name: "test", value: -1.5},
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			attr:    attr{name: "test", value: 0},
			want:    Gauge{name: "test", value: 0},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauge := Gauge{name: tc.attr.name}
			err := gauge.SetValue(tc.attr.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tc.want, gauge)
		})
	}
}
