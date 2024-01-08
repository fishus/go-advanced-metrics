package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type ValueMetricsHandlerSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *ValueMetricsHandlerSuite) SetupSuite() {
	s.ts = httptest.NewServer(ServerRouter())
	s.client = resty.New().SetBaseURL(s.ts.URL)
}

func (s *ValueMetricsHandlerSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *ValueMetricsHandlerSuite) SetupSubTest() {
	storage = metrics.NewMemStorage()
	_ = storage.AddCounter("a", 5)
	_ = storage.SetGauge("a", 1.5)
}

func (s *ValueMetricsHandlerSuite) requestValue(data []byte) *resty.Response {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(data).
		Post("value/")
	s.Require().NoError(err)

	return resp
}

func (s *ValueMetricsHandlerSuite) TestValueMetricsHandler() {
	testCases := []struct {
		name   string
		input  string
		want   string
		status int
	}{
		{
			name:   "Positive case: Counter",
			input:  `{"id":"a", "type":"counter"}`,
			want:   `{"id":"a", "type":"counter", "delta":5}`,
			status: http.StatusOK,
		},
		{
			name:   "Negative case: Counter not found",
			input:  `{"id":"b", "type":"counter"}`,
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Positive case: Gauge",
			input:  `{"id":"a", "type":"gauge"}`,
			want:   `{"id":"a", "type":"gauge", "value":1.5}`,
			status: http.StatusOK,
		},
		{
			name:   "Negative case: Gauge not found",
			input:  `{"id":"b", "type":"gauge"}`,
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Metric name not specified",
			input:  `{"id":"", "type":"counter"}`,
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Metric type not specified",
			input:  `{"id":"a", "type":""}`,
			want:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect metric type",
			input:  `{"id":"a", "type":"histogram"}`,
			want:   "",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp := s.requestValue([]byte(tc.input))
			s.Equal(tc.status, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				s.NotEmpty(resp.Body())
				s.JSONEq(tc.want, string(resp.Body()))
			}
		})
	}
}

func TestValueMetricsHandlerSuite(t *testing.T) {
	suite.Run(t, new(ValueMetricsHandlerSuite))
}
