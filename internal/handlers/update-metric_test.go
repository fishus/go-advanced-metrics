package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUpdateMetricRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func TestUpdateMetricHandler(t *testing.T) {
	ts := httptest.NewServer(ServerRouter())
	defer ts.Close()

	testCases := []struct {
		name   string
		method string
		url    string
		status int
	}{
		{
			name:   "Positive case: Counter",
			method: http.MethodPost,
			url:    "/update/counter/someMetric/123",
			status: http.StatusOK,
		},
		{
			name:   "Positive case: Gauge",
			method: http.MethodPost,
			url:    "/update/gauge/someMetric/12.34",
			status: http.StatusOK,
		},
		{
			name:   "Negative case: Method Get",
			method: http.MethodGet,
			url:    "/update/counter/someMetric/123",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "Negative case: Method Put",
			method: http.MethodPut,
			url:    "/update/counter/someMetric/123",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "Negative case: Wrong url #1",
			method: http.MethodPost,
			url:    "/update",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #3",
			method: http.MethodPost,
			url:    "/update/counter",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #4",
			method: http.MethodPost,
			url:    "/update/counter/",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #5",
			method: http.MethodPost,
			url:    "/update/counter/someMetric",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #6",
			method: http.MethodPost,
			url:    "/update/counter//",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Metric type not specified",
			method: http.MethodPost,
			url:    "/update//someMetric/1",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Metric name not specified",
			method: http.MethodPost,
			url:    "/update/counter//1",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Wrong url #7",
			method: http.MethodPost,
			url:    "/update///1",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Wrong url #8",
			method: http.MethodPost,
			url:    "/update///",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Incorrect metric type",
			method: http.MethodPost,
			url:    "/update/histogram/someMetric/1",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Counter value not specified",
			method: http.MethodPost,
			url:    "/update/counter/someMetric/",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Incorrect counter value #1",
			method: http.MethodPost,
			url:    "/update/counter/someMetric/-1",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect counter value #2",
			method: http.MethodPost,
			url:    "/update/counter/someMetric/12.34",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Incorrect counter value #3",
			method: http.MethodPost,
			url:    "/update/counter/someMetric/none",
			status: http.StatusBadRequest,
		},
		{
			name:   "Negative case: Gauge value not specified #1",
			method: http.MethodPost,
			url:    "/update/gauge/someMetric",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Gauge value not specified #2",
			method: http.MethodPost,
			url:    "/update/gauge/someMetric/",
			status: http.StatusNotFound,
		},
		{
			name:   "Negative case: Incorrect gauge value",
			method: http.MethodPost,
			url:    "/update/gauge/someMetric/none",
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := testUpdateMetricRequest(t, ts, tc.method, tc.url)
			defer resp.Body.Close()
			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}
