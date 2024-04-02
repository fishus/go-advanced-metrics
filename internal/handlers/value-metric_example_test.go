package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func ExampleValueMetricHandler_gauge() {
	storage := store.NewMemStorage()
	handlers.SetStorage(storage)

	_ = storage.SetGauge("a", 1.23)

	metricType := "gauge"
	metricID := "a"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%v/%v", metricType, metricID), nil)

	chiCtx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("metricType", metricType)
	chiCtx.URLParams.Add("metricID", metricID)

	w := httptest.NewRecorder()
	handlers.ValueMetricHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))

	// Output:
	// 200
	// 1.23
}

func ExampleValueMetricHandler_counter() {
	storage := store.NewMemStorage()
	handlers.SetStorage(storage)

	_ = storage.AddCounter("a", 123)

	metricType := "counter"
	metricID := "a"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%v/%v", metricType, metricID), nil)

	chiCtx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("metricType", metricType)
	chiCtx.URLParams.Add("metricID", metricID)

	w := httptest.NewRecorder()
	handlers.ValueMetricHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))

	// Output:
	// 200
	// 123
}
