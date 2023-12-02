package main

import (
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, handlers.UpdateHandler)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
