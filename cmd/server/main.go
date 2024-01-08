package main

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

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
	logger.Log.Info("Running server", zap.String("address", config.serverAddr), zap.String("event", "start server"))
	err := http.ListenAndServe(config.serverAddr, handlers.ServerRouter())
	if err != nil {
		logger.Log.Panic(err.Error(), zap.String("event", "start server"))
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
			logger.Log.Warn(err.Error(), zap.String("event", "load metrics from file"))
			return
		}
		logger.Log.Debug("Metric values loaded from file", zap.String("event", "load metrics from file"))
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
					logger.Log.Error(err.Error(), zap.String("event", "save metrics into file"))
				} else {
					logger.Log.Debug("Metric values saved into file", zap.String("event", "save metrics into file"))
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
