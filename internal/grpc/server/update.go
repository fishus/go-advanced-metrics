package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sg "github.com/fishus/go-advanced-metrics/internal/grpc"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

func (s *MetricsServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var response pb.UpdateResponse

	metric, err := sg.ProtoToMetric(in.Metric)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metric, code, err := Controller.UpdateMetrics(ctx, metric)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metric", metric))
		return nil, status.Error(sg.HTTPCodeToGRPC(code), err.Error())
	}

	m, err := sg.MetricToProto(metric)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metric", metric))
		return nil, status.Error(sg.HTTPCodeToGRPC(code), err.Error())

	}
	response.Metric = &m
	return &response, nil
}
