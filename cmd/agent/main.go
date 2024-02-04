package main

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/agent"
	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

func main() {
	_ = agent.Initialize()
	if err := logger.Initialize(agent.Config.LogLevel()); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)
	agent.CollectAndPostMetrics(ctx)
}
