package main

import (
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"log"
	"net/http"
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
