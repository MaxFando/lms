package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func createMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, mock
}

func TestCreateDraw(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	startTime := time.Date(2025, time.May, 5, 10, 0, 0, 0, moscowLocation)
	endTime := startTime.Add(1 * time.Hour)

	draw := &entity.Draw{
		LotteryType: "5 from 36",
		StartTime:   startTime,
		EndTime:     endTime,
	}

	mock.ExpectQuery("INSERT INTO draw_service.draw \\(lottery_type, start_time, end_time, status\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, lottery_type, start_time, end_time, status").
		WithArgs(draw.LotteryType, draw.StartTime, draw.EndTime, "PLANNED").
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(1, draw.LotteryType, draw.StartTime, draw.EndTime, "PLANNED"))

	createdDraw, err := repo.CreateDraw(context.Background(), draw)
	assert.NoError(t, err)
	assert.NotNil(t, createdDraw)
	assert.Equal(t, int32(1), createdDraw.ID)
	assert.Equal(t, entity.LotteryType5from36, createdDraw.LotteryType)

	assert.True(t, createdDraw.StartTime.In(moscowLocation).Equal(startTime))
	assert.True(t, createdDraw.EndTime.In(moscowLocation).Equal(endTime))
}

func TestGetActiveDraws(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	draws := []*entity.Draw{
		{ID: 1, LotteryType: "5 from 36", StartTime: time.Date(2025, time.May, 5, 10, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(1 * time.Hour), Status: "ACTIVE"},
		{ID: 2, LotteryType: "6 from 49", StartTime: time.Date(2025, time.May, 5, 11, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(1 * time.Hour), Status: "ACTIVE"},
	}

	mock.ExpectQuery(`SELECT id, lottery_type, start_time, end_time, status FROM draw_service.draw WHERE status = \$1`).
		WithArgs("ACTIVE").
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(draws[0].ID, draws[0].LotteryType, draws[0].StartTime, draws[0].EndTime, draws[0].Status).
			AddRow(draws[1].ID, draws[1].LotteryType, draws[1].StartTime, draws[1].EndTime, draws[1].Status))

	activeDraws, err := repo.GetActiveDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, activeDraws, 2)
	assert.Equal(t, draws[0], activeDraws[0])
	assert.Equal(t, draws[1], activeDraws[1])

	assert.True(t, activeDraws[0].StartTime.In(moscowLocation).Equal(draws[0].StartTime))
	assert.True(t, activeDraws[1].StartTime.In(moscowLocation).Equal(draws[1].StartTime))
}

func TestCancelDraw(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	drawID := int64(1)

	mock.ExpectExec(`^UPDATE draw_service.draw SET status = \$1 WHERE id = \$2`).
		WithArgs("CANCELLED", drawID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CancelDraw(context.Background(), drawID)
	assert.NoError(t, err)
}
