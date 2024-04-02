package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func ExampleListHandler() {
	storage := store.NewMemStorage()
	handlers.SetStorage(storage)

	_ = storage.SetGauge("a", 1.23)
	_ = storage.SetGauge("b", 3.21)

	_ = storage.AddCounter("a", 123)
	_ = storage.AddCounter("b", 321)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	handlers.ListHandler(w, req)
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
	// <h3>Counters:</h3>
	// <ul data-id="counters">
	// <li><strong>a</strong>: <span>123</span></li>
	// <li><strong>b</strong>: <span>321</span></li>
	// </ul>
	// <h3>Gauges:</h3>
	// <ul data-id="gauges">
	// <li><strong>a</strong>: <span>1.23</span></li>
	// <li><strong>b</strong>: <span>3.21</span></li>
	// </ul>
	//
}
