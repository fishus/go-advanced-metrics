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

func ExampleValueMetricsHandler_gauge() {
	storage := store.NewMemStorage()
	_ = handlers.NewServer(handlers.Config{
		Storage: storage,
	})

	_ = storage.SetGauge("a", 20.30)

	data := `{"id":"a", "type":"gauge"}`

	req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	handlers.ValueMetricsHandler(w, req)
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
	// {"value":20.3,"id":"a","type":"gauge"}
}

func ExampleValueMetricsHandler_counter() {
	storage := store.NewMemStorage()
	_ = handlers.NewServer(handlers.Config{
		Storage: storage,
	})

	_ = storage.AddCounter("a", 10)

	data := `{"id":"a", "type":"counter"}`

	req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	handlers.ValueMetricsHandler(w, req)
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
	// {"delta":10,"id":"a","type":"counter"}
}
