package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/controller"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

type UpdatesMetricsHandlerSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *UpdatesMetricsHandlerSuite) SetupSuite() {
	s.ts = httptest.NewServer(ServerRouter())
	s.client = resty.New().SetBaseURL(s.ts.URL)
}

func (s *UpdatesMetricsHandlerSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *UpdatesMetricsHandlerSuite) SetupSubTest() {
	config.Storage = store.NewMemStorage()
	_ = config.Storage.AddCounter("a", 7)
	_ = config.Storage.SetGauge("a", 11.15)
	controller.Storage = config.Storage
}

func (s *UpdatesMetricsHandlerSuite) requestUpdate(data []byte) *resty.Response {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(data).
		Post("updates/")
	s.Require().NoError(err)

	return resp
}

func (s *UpdatesMetricsHandlerSuite) TestUpdatesMetricsHandler() {
	testCases := []struct {
		name   string
		input  string
		want   string
		status int
	}{
		{
			name:   "Positive case: Counters",
			input:  `[{"id":"a", "type":"counter", "delta":19},{"id":"b", "type":"counter", "delta":21}]`,
			want:   `[{"id":"a", "type":"counter", "delta":26},{"id":"b", "type":"counter", "delta":21}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Double Counters",
			input:  `[{"id":"a", "type":"counter", "delta":19},{"id":"a", "type":"counter", "delta":21}]`,
			want:   `[{"id":"a", "type":"counter", "delta":47}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Gauges",
			input:  `[{"id":"a", "type":"gauge", "value":12.34},{"id":"b", "type":"gauge", "value":43.21}]`,
			want:   `[{"id":"a", "type":"gauge", "value":12.34},{"id":"b", "type":"gauge", "value":43.21}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Double Gauges",
			input:  `[{"id":"a", "type":"gauge", "value":12.34},{"id":"a", "type":"gauge", "value":43.21}]`,
			want:   `[{"id":"a", "type":"gauge", "value":43.21}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Counters and Gauges",
			input:  `[{"id":"a", "type":"counter", "delta":19},{"id":"b", "type":"counter", "delta":21},{"id":"a", "type":"gauge", "value":12.34},{"id":"b", "type":"gauge", "value":43.21}]`,
			want:   `[{"id":"a", "type":"counter", "delta":26},{"id":"b", "type":"counter", "delta":21},{"id":"a", "type":"gauge", "value":12.34},{"id":"b", "type":"gauge", "value":43.21}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Counter with Gauge value",
			input:  `[{"id":"a", "type":"counter", "delta":11, "value":12.34}]`,
			want:   `[{"id":"a", "type":"counter", "delta":18}]`,
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Gauge with Counter delta",
			input:  `[{"id":"a", "type":"gauge", "value":12.34, "delta":11}]`,
			want:   `[{"id":"a", "type":"gauge", "value":12.34}]`,
			status: http.StatusOK,
		},
		{
			name:   "Negative case: Metric name not specified",
			input:  `[{"id":"", "type":"counter", "delta":11, "value":12.34}]`,
			want:   "[]",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Metric type not specified",
			input:  `[{"id":"a", "type":"", "delta":11, "value":12.34}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect metric type",
			input:  `[{"id":"a", "type":"histogram", "delta":11, "value":12.34}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Counter value not specified",
			input:  `[{"id":"a", "type":"counter"}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect counter value #1",
			input:  `[{"id":"a", "type":"counter", "delta":-1}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect counter value #2",
			input:  `[{"id":"a", "type":"counter", "delta":12.34}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect counter value #3",
			input:  `[{"id":"a", "type":"counter", "delta":none}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Gauge value not specified",
			input:  `[{"id":"a", "type":"gauge"}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect gauge value",
			input:  `[{"id":"a", "type":"gauge", "value":none}]`,
			want:   "[]",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp := s.requestUpdate([]byte(tc.input))
			s.Equal(tc.status, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				s.NotEmpty(resp.Body())

				s.JSONEq(tc.want, string(resp.Body()))
			}
		})
	}
}

func TestUpdatesMetricsHandlerSuite(t *testing.T) {
	suite.Run(t, new(UpdatesMetricsHandlerSuite))
}
