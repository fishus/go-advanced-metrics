package main

import (
	"context"
	"errors"
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
	dbPool, _ := db.Pool()
	if dbPool != nil {
		dbStorage := storage.NewDBStorage(dbPool)
		dbStorage.MigrateCreateSchema(context.Background())
		handlers.SetStorage(dbStorage)
		return
	}

	if config.fileStoragePath != "" {
		s := storage.NewFileStorage(config.fileStoragePath)
		if config.storeInterval == 0 {
			s.SetIsSyncSave(true)
		}
		handlers.SetStorage(s)
		return
	}

	handlers.SetStorage(storage.NewMemStorage())
}

func runServer() {
	handlers.SetSecretKey(config.secretKey)
	server := &http.Server{Addr: config.serverAddr, Handler: handlers.ServerRouter()}

	go saveMetricsOnExit(server)

	logger.Log.Info("Running server", logger.String("address", config.serverAddr), logger.String("event", "start server"))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error(err.Error(), logger.String("event", "start server"))
	}
}

func loadMetricsFromFile() {
	if !config.isReqRestore {
		return
	}

	s, ok := handlers.Storage().(storage.Loader)
	if !ok {
		return
	}

	err := s.Load()
	if err != nil {
		logger.Log.Warn(err.Error(), logger.String("event", "load metrics from file"))
		return
	}
}

func saveMetricsAtIntervals() {
	if config.storeInterval <= 0 {
		return
	}

	s, ok := handlers.Storage().(storage.Saver)
	if !ok {
		return
	}

	now := time.Now()
	storeAfter := now.Add(config.storeInterval)
	for {
		now = time.Now()
		if now.After(storeAfter) {
			storeAfter = now.Add(config.storeInterval)
			err := s.Save()
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
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

	if s, ok := handlers.Storage().(storage.Saver); ok {
		err := s.Save()
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
		}
	}

	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "stop server"))
	}
}
