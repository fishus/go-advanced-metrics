package main

import (
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"net/http"
)

func main() {
	parseFlags()
	fmt.Println("Running server on", serverAddr)
	err := http.ListenAndServe(serverAddr, handlers.ServerRouter())
	if err != nil {
		panic(err)
	}
}
