package agent

import (
	"context"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type IAgentClient interface {
	Init() error
	RetryUpdateBatch(context.Context, []metrics.Metrics) error
}
