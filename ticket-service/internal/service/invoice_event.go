package service

import (
	"context"
	"encoding/json"
	"github.com/MaxFando/lms/ticket-service/internal/usecase"
	"github.com/redis/go-redis/v9"
)

type InvoiceEventHandler struct {
	redisClient   *redis.Client
	ticketUsecase *usecase.TicketUsecase
	channel       string
}

func NewInvoiceEventHandler(rdb *redis.Client, uc *usecase.TicketUsecase, channel string) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		redisClient:   rdb,
		ticketUsecase: uc,
		channel:       channel,
	}
}

func (h *InvoiceEventHandler) Run(ctx context.Context) error {
	pubsub := h.redisClient.Subscribe(ctx, h.channel)
	defer pubsub.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-pubsub.Channel():
			var ev struct {
				Type     string `json:"type"`
				TicketID int32  `json:"ticket_id"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				// todo log
				continue
			}
			if ev.Type == "invoice_overdue" || ev.Type == "invoice_failure" {
				if err := h.ticketUsecase.ReleaseBooking(ctx, ev.TicketID); err != nil {
					// todo log
				}
			} else {
				// todo log
			}
		}
	}
}
