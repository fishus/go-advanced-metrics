package handlers

import (
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testValueRequest(t *testing.T, client *resty.Client, method, url string) *resty.Response {
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		Execute(method, url)
	require.NoError(t, err)

	return resp
}

func runTestServer() *httptest.Server {
	return httptest.NewServer(ServerRouter())
}

func TestValueHandler(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	client := resty.New()

	testCases := []struct {
		name   string
		setUp  func()
		method string
		url    string
		want   string
		status int
	}{
		{
			name: "Positive case: Counter",
			setUp: func() {
				storage = metrics.NewMemStorage()
				_ = storage.AddCounter("a", 5)
			},
			method: http.MethodGet,
			url:    "/value/counter/a",
			want:   "5",
			status: http.StatusOK,
		},
		{
			name: "Positive case: Gauge",
			setUp: func() {
				storage = metrics.NewMemStorage()
				_ = storage.SetGauge("a", 1.5)
			},
			method: http.MethodGet,
			url:    "/value/gauge/a",
			want:   "1.5",
			status: http.StatusOK,
		},
		{
			name: "Positive case: Counter not found",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value/counter/a",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Positive case: Gauge not found",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value/gauge/a",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Incorrect metric type",
			setUp: func() {
				storage = metrics.NewMemStorage()
				_ = storage.AddCounter("a", 5)
				_ = storage.SetGauge("b", 1.5)
			},
			method: http.MethodGet,
			url:    "/value/histogram/h",
			want:   "",
			status: http.StatusBadRequest,
		},
		{
			name: "Negative case: Wrong url #1",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value/counter/",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Wrong url #2",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value/counter",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Wrong url #3",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value/",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Wrong url #4",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Wrong url #5",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value//",
			want:   "",
			status: http.StatusNotFound,
		},
		{
			name: "Negative case: Metric type not specified",
			setUp: func() {
				storage = metrics.NewMemStorage()
			},
			method: http.MethodGet,
			url:    "/value//name",
			want:   "",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setUp()

			resp := testValueRequest(t, client, tc.method, ts.URL+tc.url)

			assert.Equal(t, tc.status, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				assert.Equal(t, tc.want, string(resp.Body()))
			}
		})
	}
}
