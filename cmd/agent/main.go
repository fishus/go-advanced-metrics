package main

import (
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	collectAndSendMetrics()
}
