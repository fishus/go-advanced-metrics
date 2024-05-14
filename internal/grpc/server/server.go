package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

type server struct {
	server *grpc.Server
}

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func NewServer(cfg Config) *server {
	config = cfg

	srv := &server{}
	srv.server = grpc.NewServer()
	pb.RegisterMetricsServer(srv.server, &MetricsServer{})
	return srv
}

func (s *server) Run(ctx context.Context) {

	listen, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		logger.Log.Panic(err.Error())
	}

	logger.Log.Info("Running gRPC server", logger.String("address", config.ServerAddr), logger.String("event", "start server"))

	if err := s.server.Serve(listen); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}
