package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MaxFando/lms/payment-service/internal/entity"
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

func (p *Publisher) Close() error {
	return p.client.Close()
}

type event struct {
	Type     entity.EventType `json:"type"`
	TicketID int64            `json:"ticket_id"`
}

func (p *Publisher) PublishInvoice(ctx context.Context, invoice *entity.Invoice, eventType entity.EventType) error {
	data, err := json.Marshal(event{
		Type:     eventType,
		TicketID: invoice.Ticket.ID,
	})

	if err != nil {
		return err
	}

	err = p.client.Publish(ctx, p.channel, data).Err()
	if err != nil {
		return err
	}

	return nil
}
