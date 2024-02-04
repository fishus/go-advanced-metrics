package main

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/agent"
	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var config agent.Config

func main() {
	config = agent.LoadConfig()
	if err := logger.Initialize(config.LogLevel()); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)
	collectAndPostMetrics(ctx)
}
