package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGauge_Set(t *testing.T) {
	testCases := []struct {
		name    string
		gauge   Gauge
		value   float64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			value:   2.2,
			gauge:   Gauge(1.0),
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			value:   -1.5,
			gauge:   Gauge(1.0),
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			value:   0,
			gauge:   Gauge(1.0),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gauge.Set(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.value, float64(tc.gauge))
		})
	}
}
