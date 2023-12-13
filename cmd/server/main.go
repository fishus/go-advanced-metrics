package main

import (
	"log"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
)

func main() {
	parseFlags()
	runServer()
}

func runServer() {
	log.Println("Running server on", serverAddr)
	err := http.ListenAndServe(serverAddr, handlers.ServerRouter())
	if err != nil {
		panic(err)
	}
}
