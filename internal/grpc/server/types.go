package server

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

func protoToMetric(m *pb.Metric) (metrics.Metrics, error) {
	var mtype string
	switch m.Mtype {
	case pb.Mtype_TYPE_GAUGE:
		mtype = metrics.TypeGauge
	case pb.Mtype_TYPE_COUNTER:
		mtype = metrics.TypeCounter
	default:
		return metrics.Metrics{}, fmt.Errorf("unknown metric type: %s", m.Mtype)
	}

	return metrics.Metrics{
		Delta: m.Delta,
		Value: m.Value,
		ID:    m.Id,
		MType: mtype,
	}, nil
}

func metricToProto(metric metrics.Metrics) (pb.Metric, error) {
	var mtype pb.Mtype
	switch metric.MType {
	case metrics.TypeGauge:
		mtype = pb.Mtype_TYPE_GAUGE
	case metrics.TypeCounter:
		mtype = pb.Mtype_TYPE_COUNTER
	default:
		mtype = pb.Mtype_TYPE_UNSPECIFIED
	}

	return pb.Metric{
		Id:    metric.ID,
		Mtype: mtype,
		Delta: metric.Delta,
		Value: metric.Value,
	}, nil
}

func httpCodeToGRPC(code int) codes.Code {
	switch code {
	case http.StatusOK:
		return codes.OK
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusMethodNotAllowed:
		return codes.Unavailable
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusInternalServerError:
		return codes.Internal
	}
	return codes.Unknown
}
