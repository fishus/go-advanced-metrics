package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

func (s *MetricsServer) Updates(ctx context.Context, in *pb.UpdatesRequest) (*pb.UpdatesResponse, error) {
	var response pb.UpdatesResponse
	var metricsBatch []metrics.Metrics

	for _, metric := range in.Metrics {
		m, err := protoToMetric(metric)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		metricsBatch = append(metricsBatch, m)
	}

	fmt.Printf("BEFORE\n")
	for _, m := range metricsBatch {
		if m.MType == "counter" {
			fmt.Printf("%s: %#v\n", m.ID, *m.Delta)
		} else {
			fmt.Printf("%s: %#v\n", m.ID, *m.Value)
		}
	}

	metricsBatch, code, err := Controller.UpdatesMetrics(ctx, metricsBatch)
	if err != nil {
		logger.Log.Debug(err.Error(), logger.Any("metrics", metricsBatch))
		return nil, status.Error(httpCodeToGRPC(code), err.Error())
	}

	fmt.Printf("\nAFTER\n")
	for _, m := range metricsBatch {
		if m.MType == "counter" {
			fmt.Printf("%s: %#v\n", m.ID, *m.Delta)
		} else {
			fmt.Printf("%s: %#v\n", m.ID, *m.Value)
		}
	}

	var mb []*pb.Metric
	for _, metric := range metricsBatch {
		m, err := metricToProto(metric)
		if err != nil {
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return nil, status.Error(httpCodeToGRPC(code), err.Error())
		}
		mb = append(mb, &m)
	}

	response.Metrics = mb
	return &response, nil
}
