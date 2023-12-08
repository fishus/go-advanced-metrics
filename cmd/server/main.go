package main

import (
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"net/http"
)

const serverHost = "localhost:8080"

func main() {
	err := http.ListenAndServe(serverHost, handlers.ServerRouter())
	if err != nil {
		panic(err)
	}
}
