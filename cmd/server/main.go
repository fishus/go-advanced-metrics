package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	runServer()
}

func runServer() {
	logger.Log.Info("Running server", zap.String("address", config.serverAddr), zap.String("event", "start server"))
	err := http.ListenAndServe(config.serverAddr, handlers.ServerRouter())
	if err != nil {
		logger.Log.Panic(err.Error(), zap.String("event", "start server"))
	}
}
