package main

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)
	collectAndPostMetrics(ctx)
}
