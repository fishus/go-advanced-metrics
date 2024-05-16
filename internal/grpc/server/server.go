package server

import (
	"context"
	"google.golang.org/grpc"
	"net"

	ic "github.com/fishus/go-advanced-metrics/internal/grpc/interceptors"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	pb "github.com/fishus/go-advanced-metrics/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type server struct {
	server *grpc.Server
}

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func NewServer(cfg Config) *server {
	config = cfg

	interceptors := make([]grpc.ServerOption, 0)

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
	interceptors = append(interceptors, grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(ic.InterceptorLogger(logger.Log), loggerOpts...),
	))

	srv := &server{}
	srv.server = grpc.NewServer(interceptors...)
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
