package sqlext

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type config struct {
	maxOpenConns int
	maxIdleConns int

	connLifeTime time.Duration
	connIdleTime time.Duration

	tracerProvider trace.TracerProvider
}

func newConfig(opts ...ConnOption) (*config, error) {
	cfg := &config{
		maxOpenConns: 2,
		maxIdleConns: 2,

		connLifeTime: 60 * time.Minute,
		connIdleTime: 30 * time.Minute,
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("не удалось применить параметр подключения: %w", err)
		}
	}

	return cfg, nil
}
