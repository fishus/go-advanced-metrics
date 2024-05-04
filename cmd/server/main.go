package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/fishus/go-advanced-metrics/internal/app"
	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/storage"
)

var buildVersion string
var buildDate string
var buildCommit string

var config Config
var wg sync.WaitGroup

var server http.Server

func main() {
	app.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	config = c

	if err := logger.Initialize(config.logLevel); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	ctx, cancel := app.RegShutdown()
	defer cancel()

	if config.databaseDSN != "" {
		ctxDB, cancelDB := context.WithTimeout(ctx, (3 * time.Second))
		defer cancelDB()
		dbPool := db.Open(ctxDB, config.databaseDSN)
		defer dbPool.Close()
	}

	setStorage(ctx)
	loadMetricsFromFile(ctx)
	saveMetricsAtIntervals(ctx)
	saveMetricsOnExit(ctx)
	runServer(ctx)

	<-ctx.Done()
	Shutdown()
}

func readPrivateKey() []byte {
	if config.privateKeyPath != "" {
		privateKey, err := cryptokey.ReadKeyFile(config.privateKeyPath)
		if err == nil {
			return privateKey
		}
	}
	return nil
}

func setStorage(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

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

func runServer(ctx context.Context) {
	go func() {
		handlers.SetSecretKey(config.secretKey)
		handlers.SetPrivateKey(readPrivateKey())
		server = http.Server{Addr: config.serverAddr, Handler: handlers.ServerRouter()}

		logger.Log.Info("Running server", logger.String("address", config.serverAddr), logger.String("event", "start server"))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error(err.Error(), logger.String("event", "start server"))
		}
	}()
}

func loadMetricsFromFile(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

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

func saveMetricsAtIntervals(ctx context.Context) {
	if config.storeInterval <= 0 {
		return
	}

	s, ok := handlers.Storage().(storage.Saver)
	if !ok {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(config.storeInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := s.Save()
				if err != nil {
					logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
				}
			}
		}
	}()
}

func saveMetricsOnExit(ctx context.Context) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if s, ok := handlers.Storage().(storage.Saver); ok {
			err := s.Save()
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
			}
		}
	}()
}

func Shutdown() {
	wg.Wait()

	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "stop server"))
	}
}
