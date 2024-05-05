package app

import (
	"context"
	"os/signal"
	"syscall"
)

// Make context for graceful shutdown
func RegShutdown() (ctx context.Context, stop context.CancelFunc) {
	ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	return
}
