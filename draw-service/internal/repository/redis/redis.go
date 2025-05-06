package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client  *redis.Client
	channel string
}

func NewPublisher(connString, channel string) *Publisher {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	return &Publisher{
		client:  rdb,
		channel: channel,
	}
}

func (p *Publisher) PublishDraw(ctx context.Context, draw *entity.Draw) error {
	data, err := json.Marshal(draw)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = p.client.Publish(ctx, p.channel, data).Err()
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
