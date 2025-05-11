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

func NewPublisher(connString, channel string) (*Publisher, error) {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	rdb := redis.NewClient(opt)
	return &Publisher{
		client:  rdb,
		channel: channel,
	}, nil
}

func (p *Publisher) PublishDraw(ctx context.Context, draw *entity.Draw, eventType entity.EventType) error {
	type event struct {
		Type entity.EventType `json:"type"`
		Draw *entity.Draw     `json:"draw"`
	}

	data, err := json.Marshal(event{
		Type: eventType,
		Draw: draw,
	})

	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = p.client.Publish(ctx, p.channel, data).Err()
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
