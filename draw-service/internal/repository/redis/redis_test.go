package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_PublishDraw(t *testing.T) {
	db, mock := redismock.NewClientMock()
	channel := "draws_channel"

	publisher := &Publisher{
		client:  db,
		channel: channel,
	}

	ctx := context.Background()

	draw := &entity.Draw{
		ID:          1,
		LotteryType: "5 from 36",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		Status:      entity.StatusActive,
	}

	data, err := json.Marshal(draw)
	assert.NoError(t, err)

	mock.ExpectPublish(channel, data).SetVal(1)

	err = publisher.PublishDraw(ctx, draw, entity.EventTypeDrawActivated)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
