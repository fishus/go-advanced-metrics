package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

func TestValidateInputMetric(t *testing.T) {
	type want struct {
		httpCode int
		wantErr  bool
	}

	testCases := []struct {
		name   string
		metric metrics.Metrics
		want   want
	}{
		{
			name:   "Positive case: counter",
			metric: metrics.NewCounterMetric("PollCount").SetDelta(2),
			want: want{
				httpCode: 0,
				wantErr:  false,
			},
		},
		{
			name:   "Positive case: gauge",
			metric: metrics.NewGaugeMetric("RandomValue").SetValue(1.23),
			want: want{
				httpCode: 0,
				wantErr:  false,
			},
		},
		{
			name:   "Negative case: ID",
			metric: metrics.NewGaugeMetric("").SetValue(1.23),
			want: want{
				httpCode: http.StatusNotFound,
				wantErr:  true,
			},
		},
		{
			name: "Negative case: empty type",
			metric: metrics.Metrics{
				ID:    "PollCount",
				MType: "",
				Delta: func() *int64 {
					v := new(int64)
					*v = 2
					return v
				}(),
			},
			want: want{
				httpCode: http.StatusBadRequest,
				wantErr:  true,
			},
		},
		{
			name: "Negative case: wrong type",
			metric: metrics.Metrics{
				ID:    "PollCount",
				MType: "histogram",
				Delta: func() *int64 {
					v := new(int64)
					*v = 2
					return v
				}(),
			},
			want: want{
				httpCode: http.StatusBadRequest,
				wantErr:  true,
			},
		},
		{
			name:   "Negative case: empty delta",
			metric: metrics.NewCounterMetric("PollCount"),
			want: want{
				httpCode: http.StatusBadRequest,
				wantErr:  true,
			},
		},
		{
			name:   "Negative case: empty value",
			metric: metrics.NewGaugeMetric("PollCount"),
			want: want{
				httpCode: http.StatusBadRequest,
				wantErr:  true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInputMetric(tc.metric)
			if tc.want.wantErr {
				assert.Error(t, err)
				var ve *ValidMetricError
				require.ErrorAs(t, err, &ve)
				assert.Equal(t, tc.want.httpCode, ve.HTTPCode)
				return
			}

			assert.NoError(t, err)
		})
	}
}
