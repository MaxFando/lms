package tracer

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	otlpTrace "go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, spanName string) (context.Context, otlpTrace.Span) {
	ctx, span := GetTraceProvider().Tracer("").Start(ctx, spanName, otlpTrace.WithSpanKind(otlpTrace.SpanKindInternal))

	return ctx, span
}

func StartSpanWithOpt(ctx context.Context, spanName string, opts ...otlpTrace.SpanStartOption) (context.Context, otlpTrace.Span) {
	ctx, span := GetTraceProvider().Tracer("").Start(ctx, spanName, opts...)

	return ctx, span
}

func MarkSpanAsError(span otlpTrace.Span, err error) {
	if err != nil && span.IsRecording() {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}
