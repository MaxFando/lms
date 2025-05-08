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

	mock.ExpectQuery("INSERT INTO draw.draw \\(lottery_type, start_time, end_time, status\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, lottery_type, start_time, end_time, status").
		WithArgs(draw.LotteryType, draw.StartTime, draw.EndTime, entity.StatusPlanned).
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(1, draw.LotteryType, draw.StartTime, draw.EndTime, entity.StatusPlanned))

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
		{ID: 1, LotteryType: "5 from 36", StartTime: time.Date(2025, time.May, 5, 10, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(1 * time.Hour), Status: entity.StatusActive},
		{ID: 2, LotteryType: "6 from 49", StartTime: time.Date(2025, time.May, 5, 11, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(1 * time.Hour), Status: entity.StatusActive},
	}

	mock.ExpectQuery(`SELECT id, lottery_type, start_time, end_time, status FROM draw.draw WHERE status = \$1`).
		WithArgs(entity.StatusActive).
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

	drawID := int32(1)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	startTime := time.Date(2025, time.May, 5, 10, 0, 0, 0, moscowLocation)
	endTime := startTime.Add(1 * time.Hour)

	mock.ExpectQuery(`UPDATE draw.draw SET status = 'CANCELLED' WHERE id = \$1 RETURNING id, lottery_type, start_time, end_time, status`).
		WithArgs(drawID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(drawID, "5 from 36", startTime, endTime, "CANCELLED"))

	draw, err := repo.CancelDraw(context.Background(), drawID)
	assert.NoError(t, err)
	assert.NotNil(t, draw)
	assert.Equal(t, drawID, draw.ID)
	assert.Equal(t, entity.LotteryType("5 from 36"), draw.LotteryType)
	assert.Equal(t, entity.DrawStatus("CANCELLED"), draw.Status)
	assert.True(t, draw.StartTime.In(moscowLocation).Equal(startTime))
	assert.True(t, draw.EndTime.In(moscowLocation).Equal(endTime))
}

func TestActivateDraws(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation).Truncate(time.Millisecond)

	activatedDraws := []*entity.Draw{
		{ID: 1, LotteryType: "5 from 36", StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour), Status: entity.StatusActive},
		{ID: 2, LotteryType: "6 from 49", StartTime: now.Add(-2 * time.Hour), EndTime: now.Add(2 * time.Hour), Status: entity.StatusActive},
	}

	mock.ExpectQuery(`UPDATE draw\.draw\s+SET status = 'ACTIVE'\s+WHERE status = 'PLANNED' AND start_time <= \$1\s+RETURNING id, lottery_type, start_time, end_time, status`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(activatedDraws[0].ID, activatedDraws[0].LotteryType, activatedDraws[0].StartTime, activatedDraws[0].EndTime, activatedDraws[0].Status).
			AddRow(activatedDraws[1].ID, activatedDraws[1].LotteryType, activatedDraws[1].StartTime, activatedDraws[1].EndTime, activatedDraws[1].Status))

	result, err := repo.ActivateDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, activatedDraws[0].ID, result[0].ID)
	assert.Equal(t, activatedDraws[1].ID, result[1].ID)

	assert.True(t, result[0].StartTime.In(moscowLocation).Equal(activatedDraws[0].StartTime))
	assert.True(t, result[1].StartTime.In(moscowLocation).Equal(activatedDraws[1].StartTime))
}

func TestCompleteDraws(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation).Truncate(time.Millisecond)

	completedDraws := []*entity.Draw{
		{ID: 1, LotteryType: "5 from 36", StartTime: now.Add(-3 * time.Hour), EndTime: now.Add(-1 * time.Hour), Status: "COMPLETED"},
		{ID: 2, LotteryType: "6 from 49", StartTime: now.Add(-4 * time.Hour), EndTime: now.Add(-2 * time.Hour), Status: "COMPLETED"},
	}

	mock.ExpectQuery(`UPDATE draw\.draw\s+SET status = 'COMPLETED'\s+WHERE status = 'ACTIVE' AND end_time <= \$1\s+RETURNING id, lottery_type, start_time, end_time, status`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(completedDraws[0].ID, completedDraws[0].LotteryType, completedDraws[0].StartTime, completedDraws[0].EndTime, completedDraws[0].Status).
			AddRow(completedDraws[1].ID, completedDraws[1].LotteryType, completedDraws[1].StartTime, completedDraws[1].EndTime, completedDraws[1].Status))

	result, err := repo.CompleteDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, completedDraws[0].ID, result[0].ID)
	assert.Equal(t, completedDraws[1].ID, result[1].ID)

	assert.True(t, result[0].EndTime.In(moscowLocation).Equal(completedDraws[0].EndTime))
	assert.True(t, result[1].EndTime.In(moscowLocation).Equal(completedDraws[1].EndTime))
}

func TestGetCompletedDraws(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	completedDraws := []*entity.Draw{
		{ID: 1, LotteryType: "5 from 36", StartTime: time.Date(2025, time.May, 5, 10, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(-1 * time.Hour), Status: "COMPLETED"},
		{ID: 2, LotteryType: "6 from 49", StartTime: time.Date(2025, time.May, 5, 11, 0, 0, 0, moscowLocation), EndTime: time.Now().Add(-2 * time.Hour), Status: "COMPLETED"},
	}

	mock.ExpectQuery(`SELECT id, lottery_type, start_time, end_time, status FROM draw.draw WHERE status = 'COMPLETED'`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "lottery_type", "start_time", "end_time", "status"}).
			AddRow(completedDraws[0].ID, completedDraws[0].LotteryType, completedDraws[0].StartTime, completedDraws[0].EndTime, completedDraws[0].Status).
			AddRow(completedDraws[1].ID, completedDraws[1].LotteryType, completedDraws[1].StartTime, completedDraws[1].EndTime, completedDraws[1].Status))

	result, err := repo.GetCompletedDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, completedDraws[0], result[0])
	assert.Equal(t, completedDraws[1], result[1])

	assert.True(t, result[0].StartTime.In(moscowLocation).Equal(completedDraws[0].StartTime))
	assert.True(t, result[1].StartTime.In(moscowLocation).Equal(completedDraws[1].StartTime))
}

func TestGetDrawResultByDrawID(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	repo := NewDrawRepository(db)

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	drawID := int32(1)
	drawResult := entity.DrawResult{
		ID:                 10,
		DrawID:             drawID,
		WinningCombination: "1,5,12,23,30",
		ResultTime:         time.Date(2025, time.May, 5, 12, 0, 0, 0, moscowLocation),
	}

	mock.ExpectQuery(`SELECT id, draw_id, winning_combination, result_time FROM draw.draw_result WHERE draw_id = \$1`).
		WithArgs(drawID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "draw_id", "winning_combination", "result_time"}).
			AddRow(drawResult.ID, drawResult.DrawID, drawResult.WinningCombination, drawResult.ResultTime))

	result, err := repo.GetDrawResult(context.Background(), drawID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, drawResult, *result)
}
