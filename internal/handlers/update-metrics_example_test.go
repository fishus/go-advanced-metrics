package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func ExampleUpdateMetricsHandler_gauge() {
	storage := store.NewMemStorage()
	_ = handlers.NewServer(handlers.Config{
		Storage: storage,
	})

	_ = storage.SetGauge("a", 1.23)

	data := `{"id":"a", "type":"gauge", "value":123.45}`

	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	handlers.UpdateMetricsHandler(w, req)
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
	// {"value":123.45,"id":"a","type":"gauge"}
}

func ExampleUpdateMetricsHandler_counter() {
	storage := store.NewMemStorage()
	_ = handlers.NewServer(handlers.Config{
		Storage: storage,
	})

	_ = storage.AddCounter("a", 1)

	data := `{"id":"a", "type":"counter", "delta":10}`

	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	handlers.UpdateMetricsHandler(w, req)
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
	// {"delta":11,"id":"a","type":"counter"}
}
