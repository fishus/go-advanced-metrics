package server

import (
	"context"
	"sync"
	"time"

	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/handlers"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

var wgServer sync.WaitGroup

type IServer interface {
	Run(ctx context.Context)
	Shutdown(ctx context.Context) error
}

var Server IServer

var Storage store.MetricsStorager

func SetStorage() {
	dbPool, _ := db.Pool()
	if dbPool != nil {
		dbStorage := store.NewDBStorage(dbPool)
		dbStorage.MigrateCreateSchema(context.Background())
		Storage = dbStorage
		return
	}

	if Config.FileStoragePath() != "" {
		s := store.NewFileStorage(Config.FileStoragePath())
		if Config.StoreInterval() == 0 {
			s.SetIsSyncSave(true)
		}
		Storage = s
		return
	}

	Storage = store.NewMemStorage()
}

func LoadMetricsFromFile() {
	if !Config.IsReqRestore() {
		return
	}

	s, ok := Storage.(store.Loader)
	if !ok {
		return
	}

	err := s.Load()
	if err != nil {
		logger.Log.Warn(err.Error(), logger.String("event", "load metrics from file"))
		return
	}
}

func SaveMetricsAtIntervals(ctx context.Context) {
	if Config.StoreInterval() <= 0 {
		return
	}

	s, ok := Storage.(store.Saver)
	if !ok {
		return
	}

	wgServer.Add(1)
	go func() {
		defer wgServer.Done()

		ticker := time.NewTicker(Config.StoreInterval())
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

func SaveMetricsOnExit(ctx context.Context) {
	wgServer.Add(1)
	go func() {
		defer wgServer.Done()
		<-ctx.Done()

		if s, ok := Storage.(store.Saver); ok {
			err := s.Save()
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("event", "save metrics into file"))
			}
		}
	}()
}

func RunServer(ctx context.Context) {
	Server = handlers.NewServer(handlers.Config{
		ServerAddr:    Config.ServerAddr(),
		Storage:       Storage,
		SecretKey:     Config.SecretKey(),
		PrivateKey:    PrivateKey,
		TrustedSubnet: Config.TrustedSubnet(),
	})

	go func() {
		Server.Run(ctx)
	}()
}
