package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Shutdown implements graceful app
func Shutdown(cancel context.CancelFunc) {
	go func() {
		termSig := make(chan os.Signal, 1)
		signal.Notify(termSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-termSig
		cancel()
	}()
}
