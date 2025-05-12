package payment

import (
	"context"
	"errors"
	"math/rand"

	"github.com/MaxFando/lms/payment-service/internal/entity"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) Pay(_ context.Context, card *entity.Card) (int64, error) {
	if card == nil {
		return 0, errors.New("card is nil")
	}

	transactionID := rand.Int63()

	if card.CVV == "123" || rand.Intn(10) < 8 {
		return transactionID, nil
	}

	return 0, errors.New("payment failure")
}

func (c *Client) Refund(_ context.Context, _ int64) error {
	if rand.Intn(10) < 8 {
		return nil
	}

	return errors.New("payment failure")
}
