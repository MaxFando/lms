package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"time"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const connectTimeout = 2 * time.Second

type ShutdownFn func(context.Context) error

func InitDefaultProvider(cfg Config) (ShutdownFn, error) {
	ctx := context.Background()

	tracerProvider, err := NewProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(jaeger.Jaeger{})

	if cfg.ErrorLogFunc != nil {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			cfg.ErrorLogFunc(ctx, err, "tracer error")
		}))
	}

	return tracerProvider.Shutdown, nil
}

func NewProvider(ctx context.Context, cfg Config) (*tracesdk.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpointURL(cfg.URL),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(connectTimeout),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании OTLP экспортера: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.AppName),
		semconv.DeploymentEnvironmentKey.String(cfg.Environment),
	)

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.TraceRatio)),
	)

	return tracerProvider, nil
}

func GetTraceProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}
