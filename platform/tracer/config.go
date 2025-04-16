package tracer

import (
	"context"
)

type ClientType string

const (
	HTTPClient ClientType = "HTTP"
)

type logFunction func(ctx context.Context, err error, message string)

type Config struct {
	URL          string
	AppName      string
	Environment  string
	TraceRatio   float64
	Client       ClientType
	ErrorLogFunc logFunction
}

type ConfigOption func(*Config)

func WithAppName(name string) ConfigOption {
	return func(cfg *Config) {
		cfg.AppName = name
	}
}

func WithEnvironment(env string) ConfigOption {
	return func(cfg *Config) {
		cfg.Environment = env
	}
}

func NewConfig(url string, opts ...ConfigOption) (Config, error) {
	cfg := Config{
		URL:        url,
		TraceRatio: 1.0,
		Client:     HTTPClient,
	}

	for _, o := range opts {
		o(&cfg)
	}

	return cfg, nil
}
