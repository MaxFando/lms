package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(ctx context.Context, msg string, keyvals ...interface{})
	Info(ctx context.Context, msg string, keyvals ...interface{})
	Error(ctx context.Context, message string, keyvals ...interface{})
	With(keyvals ...interface{}) Logger
}

func NewLogger() Logger {
	h := &ContextHandler{slog.NewJSONHandler(os.Stdout, nil)}
	l := slog.New(h)

	return &logger{
		log: l,
	}
}

type logger struct {
	log *slog.Logger
}

func (l *logger) Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	keyvals = checkKeyvals(keyvals...)

	l.log.DebugContext(ctx, msg, keyvals...)
}

func (l *logger) Info(ctx context.Context, msg string, keyvals ...interface{}) {
	keyvals = checkKeyvals(keyvals...)
	l.log.InfoContext(ctx, msg, keyvals...)
}

func (l *logger) Error(ctx context.Context, msg string, keyvals ...interface{}) {
	keyvals = checkKeyvals(keyvals...)
	l.log.ErrorContext(ctx, msg, keyvals...)
}

func (l *logger) With(keyvals ...interface{}) Logger {
	keyvals = checkKeyvals(keyvals...)
	return &logger{
		log: l.log.With(keyvals...),
	}
}

func checkKeyvals(keyvals ...interface{}) []interface{} {
	if len(keyvals)%2 != 0 {
		return keyvals[:len(keyvals)-1]
	}

	return keyvals
}

type ctxKey string

const (
	uberTraceID ctxKey = "uber-trace-id"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(uberTraceID).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(uberTraceID).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, uberTraceID, v)
	}

	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(parent, uberTraceID, v)
}
