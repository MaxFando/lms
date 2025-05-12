package ticket

import (
	"context"

	"github.com/MaxFando/lms/payment-service/internal/entity"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) BookTicket(_ context.Context, _ int64, ticketID int64) (*entity.Ticket, error) {
	return &entity.Ticket{
		ID: ticketID,
	}, nil
}
