package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		request  string
		wantCode int
	}{
		{
			name:     "Positive case: Counter",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric/123",
			wantCode: http.StatusOK,
		},
		{
			name:     "Positive case: Gauge",
			method:   http.MethodPost,
			request:  "/update/gauge/someMetric/12.34",
			wantCode: http.StatusOK,
		},
		{
			name:     "Negative case: Method Get",
			method:   http.MethodGet,
			request:  "/update/counter/someMetric/123",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Negative case: Method Put",
			method:   http.MethodPut,
			request:  "/update/counter/someMetric/123",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Negative case: Metric type not specified #1",
			method:   http.MethodPost,
			request:  "/update",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Metric type not specified #2",
			method:   http.MethodPost,
			request:  "/update/",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Incorrect metric type",
			method:   http.MethodPost,
			request:  "/update/histogram",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Metric name not specified #1",
			method:   http.MethodPost,
			request:  "/update/counter",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative case: Metric name not specified #2",
			method:   http.MethodPost,
			request:  "/update/counter/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative case: Metric name not specified #3",
			method:   http.MethodPost,
			request:  "/update/gauge",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative case: Metric name not specified #4",
			method:   http.MethodPost,
			request:  "/update/gauge/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative case: Counter value not specified #1",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Counter value not specified #2",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric/",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Incorrect counter value #1",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Incorrect counter value #2",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric/12.34",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Incorrect counter value #3",
			method:   http.MethodPost,
			request:  "/update/counter/someMetric/none",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Gauge value not specified #1",
			method:   http.MethodPost,
			request:  "/update/gauge/someMetric",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Gauge value not specified #2",
			method:   http.MethodPost,
			request:  "/update/gauge/someMetric/",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Negative case: Incorrect gauge value",
			method:   http.MethodPost,
			request:  "/update/gauge/someMetric/none",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()
			UpdateHandler(w, r)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tc.wantCode, res.StatusCode)
		})
	}
}
