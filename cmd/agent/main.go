package main

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/agent"
	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	app.PrintBuildInfo(buildVersion, buildDate, buildCommit)

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
