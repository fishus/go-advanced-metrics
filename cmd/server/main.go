package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/storage"
)

var config Config

func main() {
	config = loadConfig()
	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	if config.databaseDSN != "" {
		ctxDBTimeout, cancelDBTimeout := context.WithTimeout(context.Background(), (3 * time.Second))
		defer cancelDBTimeout()
		dbPool := db.Open(ctxDBTimeout, config.databaseDSN)
		defer dbPool.Close()
	}

	setStorage()
	loadMetricsFromFile()
	go saveMetricsAtIntervals()
	runServer()
}

func setStorage() {
	handlers.Config.IsSyncMetricsSave = false

	dbConn, err := db.Pool()
	if err == nil {
		dbStorage := storage.NewDBStorage(dbConn)
		dbStorage.MigrateCreateSchema(context.Background())
		handlers.SetStorage(dbStorage)
		return
	}

	if config.fileStoragePath != "" {
		if config.storeInterval == 0 {
			handlers.Config.IsSyncMetricsSave = true
		}
		handlers.SetStorage(storage.NewFileStorage(config.fileStoragePath))
		return
	}

	handlers.SetStorage(storage.NewMemStorage())
}

func runServer() {
	server := &http.Server{Addr: config.serverAddr, Handler: handlers.ServerRouter()}

	go saveMetricsOnExit(server)

	logger.Log.Info("Running server", logger.String("address", config.serverAddr), logger.String("event", "start server"))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error(err.Error(), logger.String("event", "start server"))
	}
}

func loadMetricsFromFile() {
	if !config.isReqRestore || fmt.Sprintf("%T", handlers.Storage()) != "*storage.FileStorage" {
		return
	}

	s := handlers.Storage().(*storage.FileStorage)

	err := s.Load()
	if !errors.Is(err, storage.ErrEmptyFilename) {
		if err != nil {
			logger.Log.Warn(err.Error(), logger.String("event", "load metrics from file"))
			return
		}
		logger.Log.Debug("Metric values loaded from file", logger.String("event", "load metrics from file"))
	}
}

func saveMetricsAtIntervals() {
	if config.storeInterval <= 0 || fmt.Sprintf("%T", handlers.Storage()) != "*storage.FileStorage" {
		return
	}

	s := handlers.Storage().(*storage.FileStorage)

	now := time.Now()
	storeAfter := now.Add(config.storeInterval)
	for {
		now = time.Now()
		if now.After(storeAfter) {
			storeAfter = now.Add(config.storeInterval)
			err := s.Save()
			if !errors.Is(err, storage.ErrEmptyFilename) {
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

	if fmt.Sprintf("%T", handlers.Storage()) == "*storage.FileStorage" {
		s := handlers.Storage().(*storage.FileStorage)
		err := s.Save()
		if !errors.Is(err, storage.ErrEmptyFilename) {
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
