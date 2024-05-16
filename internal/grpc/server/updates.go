package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sg "github.com/fishus/go-advanced-metrics/internal/grpc"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

func (s *MetricsServer) Updates(ctx context.Context, in *pb.UpdatesRequest) (*pb.UpdatesResponse, error) {
	var response pb.UpdatesResponse
	var metricsBatch []metrics.Metrics

	for _, metric := range in.Metrics {
		m, err := sg.ProtoToMetric(metric)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		metricsBatch = append(metricsBatch, m)
	}

	metricsBatch, code, err := Controller.UpdatesMetrics(ctx, metricsBatch)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metrics", metricsBatch))
		return nil, status.Error(sg.HTTPCodeToGRPC(code), err.Error())
	}

	var mb []*pb.Metric
	for _, metric := range metricsBatch {
		m, err := sg.MetricToProto(metric)
		if err != nil {
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return nil, status.Error(sg.HTTPCodeToGRPC(code), err.Error())
		}
		mb = append(mb, &m)
	}

	response.Metrics = mb
	return &response, nil
}
