package redis

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MaxFando/lms/payment-service/internal/entity"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_PublishInvoice(t *testing.T) {
	db, mock := redismock.NewClientMock()
	channel := "invoice_channel"

	publisher := &Publisher{
		client:  db,
		channel: channel,
	}

	ctx := context.Background()

	invoice := &entity.Invoice{
		Ticket: &entity.Ticket{
			ID: 111111,
		},
	}

	data, err := json.Marshal(&event{
		Type:     entity.EventTypeInvoiceOverdue,
		TicketID: 111111,
	})
	assert.NoError(t, err)

	mock.ExpectPublish(channel, data).SetVal(1)

	err = publisher.PublishInvoice(ctx, invoice, entity.EventTypeInvoiceOverdue)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
