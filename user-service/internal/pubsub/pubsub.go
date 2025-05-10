package pubsub

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"

	"go.opentelemetry.io/otel/trace"
)

type PubSub interface {
	PublishUserCreated(ctx context.Context, userID int64) error
	PublishUserLoggedIn(ctx context.Context, userID int64) error
}

type redisPubSub struct {
	client *redis.Client // Changed: go-redis v9
}

func NewRedisPubSub(rdb *redis.Client) PubSub {
	return &redisPubSub{client: rdb}
}

func (r *redisPubSub) PublishUserCreated(ctx context.Context, userID int64) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PublishUserCreated")
	defer span.End()

	channel := "user.created"
	payload := fmt.Sprintf("%d", userID)
	return r.client.Publish(ctx, channel, payload).Err()
}

func (r *redisPubSub) PublishUserLoggedIn(ctx context.Context, userID int64) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PublishUserLoggedIn")
	defer span.End()

	channel := "user.logged_in"
	payload := fmt.Sprintf("%d", userID)
	return r.client.Publish(ctx, channel, payload).Err()
}
