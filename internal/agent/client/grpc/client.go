package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/encoding/gzip"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	sg "github.com/fishus/go-advanced-metrics/internal/grpc"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	pb "github.com/fishus/go-advanced-metrics/proto"
)

type Client struct {
	config Config
	conn   *grpc.ClientConn
	client pb.MetricsClient
	ip     string
}

func NewClient(conf Config) *Client {
	client := &Client{
		config: conf,
	}
	return client
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Init() error {
	conn, err := grpc.Dial(c.config.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("address", c.config.ServerAddr), logger.String("event", "start agent worker"))
		return fmt.Errorf("can't connect to grpc server: %w", err)
	}

	c.conn = conn
	c.client = pb.NewMetricsClient(conn)
	logger.Log.Info("Running gRPC worker", logger.String("address", c.config.ServerAddr), logger.String("event", "start agent worker"))

	return nil
}
func (c *Client) RetryUpdateBatch(ctx context.Context, batch []metrics.Metrics) (err error) {
	retryDelay := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
		0,
	}

	for _, delay := range retryDelay {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = c.UpdateBatch(ctx, batch)

		if err == nil {
			return nil
		}
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded,
				codes.Unavailable:
			default:
				return fmt.Errorf(e.Message())
			}
		} else {
			return err
		}

		time.Sleep(delay)
	}

	return err
}

func (c *Client) UpdateBatch(ctx context.Context, batch []metrics.Metrics) error {
	if len(batch) == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var mb []*pb.Metric
	for _, metric := range batch {
		m, err := sg.MetricToProto(metric)
		if err != nil {
			logger.Log.Debug(err.Error(), logger.Any("metric", metric))
			return fmt.Errorf("can't convert metric to proto: %w", err)
		}
		mb = append(mb, &m)
	}

	req := &pb.UpdatesRequest{
		Metrics: mb,
	}

	logger.Log.Debug(`Send Updates request`,
		logger.String("event", "send request"),
		logger.String("addr", c.config.ServerAddr),
		logger.String("data", req.String()))

	gz := grpc.UseCompressor(gzip.Name)

	resp, err := c.client.Updates(ctx, req, gz)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Debug(`Received response from the server`,
		logger.String("event", "response received"),
		logger.String("data", resp.String()))

	return nil
}
