package service

import (
	"context"
	"encoding/json"
	"github.com/MaxFando/lms/ticket-service/internal/entity"
	"github.com/MaxFando/lms/ticket-service/internal/usecase"
	"github.com/redis/go-redis/v9"
)

type DrawEventHandler struct {
	redisClient    *redis.Client
	ticketUsecase  *usecase.TicketUsecase
	ticketsPerDraw int
	channel        string
}

func NewDrawEventHandler(rdb *redis.Client, uc *usecase.TicketUsecase, channel string, ticketsPerDraw int) *DrawEventHandler {
	return &DrawEventHandler{
		redisClient:    rdb,
		ticketUsecase:  uc,
		channel:        channel,
		ticketsPerDraw: ticketsPerDraw,
	}
}

func (h *DrawEventHandler) Run(ctx context.Context) error {
	pubsub := h.redisClient.Subscribe(ctx, h.channel)
	defer pubsub.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-pubsub.Channel():
			var ev struct {
				Type string      `json:"type"`
				Draw entity.Draw `json:"draw"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				// todo log
				continue
			}
			if ev.Type != entity.EventTypeDrawActivated {
				// todo log
				continue
			}
			err := h.ticketUsecase.GenerateTickets(ctx, ev.Draw.ID, h.ticketsPerDraw)
			if err != nil {
				// todo log
			}
		}
	}
}
