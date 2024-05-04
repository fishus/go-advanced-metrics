package main

import (
	"github.com/fishus/go-advanced-metrics/internal/agent"
	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	app.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	if err := agent.Initialize(); err != nil {
		panic(err)
	}

	if err := logger.Initialize(agent.Config.LogLevel()); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	ctx, cancel := app.RegShutdown()
	defer cancel()

	agent.CollectAndPostMetrics(ctx)

	<-ctx.Done()
	agent.Shutdown()
}
