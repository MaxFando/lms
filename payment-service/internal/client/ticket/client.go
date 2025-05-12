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

func (c *Client) BookTicket(ctx context.Context, userId int64, ticketI int64) (*entity.Ticket, error) {
	return &entity.Ticket{
		ID: 111111,
	}, nil
}
