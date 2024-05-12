package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
)

type server struct {
	server *http.Server
}

func NewServer(cfg Config) *server {
	config = cfg
	return &server{}
}

func (s *server) Run(ctx context.Context) {
	s.server = &http.Server{Addr: config.ServerAddr, Handler: ServerRouter()}
	logger.Log.Info("Running server", logger.String("address", config.ServerAddr), logger.String("event", "start server"))
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error(err.Error(), logger.String("event", "start server"))
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
