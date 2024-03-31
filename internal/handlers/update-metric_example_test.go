package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func ExampleUpdateMetricHandler_gauge() {
	storage := store.NewMemStorage()
	handlers.SetStorage(storage)

	_ = storage.SetGauge("a", 1.23)

	metricType := "gauge"
	metricID := "a"
	metricValue := "12.34"

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%v/%v/%v", metricType, metricID, metricValue), nil)

	chiCtx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("metricType", metricType)
	chiCtx.URLParams.Add("metricID", metricID)
	chiCtx.URLParams.Add("metricValue", metricValue)

	w := httptest.NewRecorder()
	handlers.UpdateMetricHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}

func ExampleUpdateMetricHandler_counter() {
	storage := store.NewMemStorage()
	handlers.SetStorage(storage)

	_ = storage.AddCounter("a", 1)

	metricType := "counter"
	metricID := "a"
	metricValue := "10"

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%v/%v/%v", metricType, metricID, metricValue), nil)

	chiCtx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("metricType", metricType)
	chiCtx.URLParams.Add("metricID", metricID)
	chiCtx.URLParams.Add("metricValue", metricValue)

	w := httptest.NewRecorder()
	handlers.UpdateMetricHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}
