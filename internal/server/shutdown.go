package server

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/logger"
)

func Shutdown() {
	wgServer.Wait()

	err := Server.Shutdown(context.Background())
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "stop server"))
	}
}
