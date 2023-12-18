package main

import (
	"log"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
)

var config Config

func main() {
	config = loadConfig()
	runServer()
}

func runServer() {
	log.Println("Running server on", config.serverAddr)
	err := http.ListenAndServe(config.serverAddr, handlers.ServerRouter())
	if err != nil {
		panic(err)
	}
}
