package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/pashagolub/pgxmock/v3"

	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
)

func ExamplePingDBHandler() {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	mock.ExpectPing()
	db.SetPool(mock)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	w := httptest.NewRecorder()
	handlers.PingDBHandler(w, req)
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
}
