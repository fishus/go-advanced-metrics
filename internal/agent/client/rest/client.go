package rest

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"time"

	"github.com/go-resty/resty/v2"

	ac "github.com/fishus/go-advanced-metrics/internal/agent/client"
	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/fishus/go-advanced-metrics/internal/secure"
)

type Client struct {
	config Config
	client *resty.Client
	gz     *gzip.Writer
}

func NewClient(conf Config) *Client {
	client := &Client{
		config: conf,
	}
	return client
}

func (c *Client) Init() error {
	c.client = resty.New().SetBaseURL("http://" + c.config.ServerAddr)
	logger.Log.Info("Running rest worker", logger.String("address", c.config.ServerAddr), logger.String("event", "start agent worker"))

	ip, err := ac.GetIP()
	if err != nil {
		logger.Log.Warn(err.Error())
		return err
	} else if ip != nil {
		c.client.SetHeader("X-Real-IP", ip.String())
	}

	gz, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		logger.Log.Warn(err.Error())
		return err
	}
	c.gz = gz
	return nil
}

func (c *Client) RetryUpdateBatch(ctx context.Context, batch []metrics.Metrics) (err error) {
	var neterr *net.OpError

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

		errors.As(err, &neterr)
		if err == nil || !errors.Is(err, neterr) {
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

	jsonBody, err := json.Marshal(batch)
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "encode json"),
			logger.Any("data", batch))
		return err
	}

	var hashString string
	if c.config.SecretKey != "" {
		hash := secure.Hash(jsonBody, []byte(c.config.SecretKey))
		hashString = hex.EncodeToString(hash[:])
	}

	if len(c.config.PublicKey) > 0 {
		jsonBody, err = cryptokey.Encrypt(jsonBody, c.config.PublicKey)
		if err != nil {
			return err
		}
	}

	buf := bytes.NewBuffer(nil)
	c.gz.Reset(buf)
	_, err = c.gz.Write(jsonBody)
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "compress request"),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}
	err = c.gz.Close()
	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "compress request"),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}

	logger.Log.Debug(`Send POST /updates/ request`,
		logger.String("event", "send request"),
		logger.String("addr", c.config.ServerAddr),
		logger.Any("body", json.RawMessage(jsonBody)))

	req := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(buf)

	if hashString != "" {
		req.SetHeader("HashSHA256", hashString)
	}

	url := "updates/"
	resp, err := req.Post(url)

	if err != nil {
		logger.Log.Error(err.Error(),
			logger.String("event", "send request"),
			logger.String("url", "http://"+c.config.ServerAddr+"/"+url),
			logger.Any("body", json.RawMessage(jsonBody)))
		return err
	}

	rawBody := resp.RawBody()
	defer rawBody.Close()

	gzBody, err := gzip.NewReader(rawBody)
	if err != nil && err != io.EOF {
		return err
	}
	defer gzBody.Close()

	body, err := io.ReadAll(gzBody)
	if err != nil && err != io.EOF {
		return err
	}

	logger.Log.Debug(`Received response from the server`, logger.String("event", "response received"), logger.Any("headers", resp.Header()), logger.Any("body", json.RawMessage(body)))

	return nil
}
