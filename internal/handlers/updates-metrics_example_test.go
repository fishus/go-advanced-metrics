package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/fishus/go-advanced-metrics/internal/controller"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

func ExampleUpdatesMetricsHandler() {
	storage := store.NewMemStorage()
	_ = handlers.NewServer(handlers.Config{
		Storage: storage,
	})

	_ = storage.SetGauge("a", 1.23)
	_ = storage.AddCounter("a", 1)
	controller.Storage = storage

	data := `
	[
		{
			"id": "a",
			"type": "counter",
			"delta": 2
		},
		{
			"id": "a",
			"type": "counter",
			"delta": 3
		},
		{
			"id": "b",
			"type": "counter",
			"delta": 5
		},
		{
			"id": "a",
			"type": "gauge",
			"value": 12.34
		},
		{
			"id": "a",
			"type": "gauge",
			"value": 23.45
		},
		{
			"id": "b",
			"type": "gauge",
			"value": 43.21
		}
	]`

	req := httptest.NewRequest(http.MethodPost, "/updates/", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	handlers.UpdatesMetricsHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	jsonString, err := PrettyJSONString(body)
	if err != nil {
		return
	}
	fmt.Println(jsonString)

	// Output:
	// 200
	// [
	//     {
	//         "delta": 6,
	//         "id": "a",
	//         "type": "counter"
	//     },
	//     {
	//         "delta": 5,
	//         "id": "b",
	//         "type": "counter"
	//     },
	//     {
	//         "value": 23.45,
	//         "id": "a",
	//         "type": "gauge"
	//     },
	//     {
	//         "value": 43.21,
	//         "id": "b",
	//         "type": "gauge"
	//     }
	// ]
}

func PrettyJSONString(data []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
