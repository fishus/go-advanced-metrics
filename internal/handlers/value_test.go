package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type ValueHandlerSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *ValueHandlerSuite) SetupSuite() {
	s.ts = httptest.NewServer(ServerRouter())
	s.client = resty.New()
}

func (s *ValueHandlerSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *ValueHandlerSuite) SetupTest() {
	storage = metrics.NewMemStorage()
	_ = storage.AddCounter("a", 5)
	_ = storage.SetGauge("a", 1.5)
}

func (s *ValueHandlerSuite) requestValue(url string) *resty.Response {
	resp, err := s.client.R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		Get(url)
	s.Require().NoError(err)
	return resp
}

func (s *ValueHandlerSuite) TestValueHandler() {
	testCases := []struct {
		name   string
		url    string
		want   string
		status int
	}{
		{
			name:   "Positive case: Counter",
			url:    "/value/counter/a",
			want:   "5",
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Gauge",
			url:    "/value/gauge/a",
			want:   "1.5",
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Counter not found",
			url:    "/value/counter/x",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Positive case: Gauge not found",
			url:    "/value/gauge/x",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Incorrect metric type",
			url:    "/value/histogram/h",
			want:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Wrong url #1",
			url:    "/value/counter/",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #2",
			url:    "/value/counter",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #3",
			url:    "/value/",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #4",
			url:    "/value",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #5",
			url:    "/value//",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Metric type not specified",
			url:    "/value//name",
			want:   "",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp := s.requestValue(s.ts.URL + tc.url)
			s.Equal(tc.status, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				s.Equal(tc.want, string(resp.Body()))
			}
		})
	}
}

func TestValueHandlerSuite(t *testing.T) {
	suite.Run(t, new(ValueHandlerSuite))
}
