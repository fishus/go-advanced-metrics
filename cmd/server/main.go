package main

import (
	"context"
	"time"

	"github.com/fishus/go-advanced-metrics/internal/app"
	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	app.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	if err := server.Initialize(); err != nil {
		panic(err)
	}

	if err := logger.Initialize(server.Config.LogLevel()); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	ctx, cancel := app.RegShutdown()
	defer cancel()

	if server.Config.DatabaseDSN() != "" {
		ctxDB, cancelDB := context.WithTimeout(ctx, (3 * time.Second))
		defer cancelDB()
		dbPool := db.Open(ctxDB, server.Config.DatabaseDSN())
		defer dbPool.Close()
	}

	server.SetStorage()
	server.LoadMetricsFromFile()
	server.SaveMetricsAtIntervals(ctx)
	server.SaveMetricsOnExit(ctx)
	server.RunServer(ctx)

	<-ctx.Done()
	server.Shutdown()
}
