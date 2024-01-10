package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	loadMetricsFromFile()
	go saveMetricsAtIntervals()
	runServer()
}

func runServer() {
	if config.storeInterval == 0 {
		handlers.Config.IsSyncMetricsSave = true
	} else {
		handlers.Config.IsSyncMetricsSave = false
	}

	st := handlers.Storage()
	st.Filename = config.fileStoragePath

	server := &http.Server{Addr: config.serverAddr, Handler: handlers.ServerRouter()}

	go saveMetricsOnExit(server)

	logger.Log.Info("Running server", logger.String("address", config.serverAddr), logger.String("event", "start server"))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error(err.Error(), logger.String("event", "start server"))
	}
}

func loadMetricsFromFile() {
	if !config.isReqRestore || config.fileStoragePath == "" {
		return
	}

	s := handlers.Storage()
	s.Filename = config.fileStoragePath

	err := s.Load()
	if !errors.Is(err, metrics.ErrEmptyFilename) {
		if err != nil {
			logger.Log.Warn(err.Error(), logger.String("event", "load metrics from file"))
			return
		}
		logger.Log.Debug("Metric values loaded from file", logger.String("event", "load metrics from file"))
	}
}

func saveMetricsAtIntervals() {
	if config.storeInterval <= 0 || config.fileStoragePath == "" {
		return
	}

	s := handlers.Storage()
	s.Filename = config.fileStoragePath

	now := time.Now()
	storeAfter := now.Add(config.storeInterval)
	for {
		now = time.Now()
		if now.After(storeAfter) {
			storeAfter = now.Add(config.storeInterval)
			err := s.Save()
			if !errors.Is(err, metrics.ErrEmptyFilename) {
				if err != nil {
					logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
				} else {
					logger.Log.Debug("Metric values saved into file", logger.String("event", "save metrics into file"))
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func saveMetricsOnExit(server *http.Server) {
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-termSig
	logger.Log.Debug("Server interrupt signal caught", logger.String("event", "stop server"), logger.String("signal", sig.String()))

	if config.fileStoragePath != "" {
		s := handlers.Storage()
		s.Filename = config.fileStoragePath
		err := s.Save()
		if !errors.Is(err, metrics.ErrEmptyFilename) {
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
			} else {
				logger.Log.Debug("Metric values saved into file", logger.String("event", "save metrics into file"))
			}
		}
	}

	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "stop server"))
	}
}
