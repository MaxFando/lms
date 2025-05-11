package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/MaxFando/lms/draw-service/tests/containers"
	"github.com/MaxFando/lms/platform/sqlext"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()

	pgC, err := containers.CreatePostgresContainer(ctx)
	require.NoError(t, err)

	db, err := sqlext.OpenSqlxViaPgxConnPool(context.Background(), pgC.ConnectionString)
	require.NoError(t, err)

	goose.SetBaseFS(nil)

	migrationsDir := "../../../migrations"
	sqlDB := db.DB

	err = goose.Up(sqlDB, migrationsDir)
	require.NoError(t, err)

	cleanup := func() {
		_ = db.Close()
		_ = pgC.Terminate(ctx)
	}

	return db, cleanup
}

func TestCreateDraw(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)

	loc, _ := time.LoadLocation("Europe/Moscow")
	start := time.Date(2025, 5, 5, 10, 0, 0, 0, loc)
	end := start.Add(time.Hour)

	draw := &entity.Draw{
		LotteryType: entity.LotteryType5from36,
		StartTime:   start,
		EndTime:     end,
	}

	created, err := repo.CreateDraw(context.Background(), draw)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, entity.LotteryType5from36, created.LotteryType)
	require.Equal(t, entity.StatusPlanned, created.Status)

	require.True(t, created.StartTime.In(loc).Equal(start))
	require.True(t, created.EndTime.In(loc).Equal(end))
}

func TestGetActiveDraws(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	draws := []*entity.Draw{
		{LotteryType: "5 from 36", StartTime: time.Now().Add(-1 * time.Hour).In(moscowLocation), EndTime: time.Now().Add(1 * time.Hour).In(moscowLocation), Status: entity.StatusActive},
		{LotteryType: "6 from 49", StartTime: time.Now().Add(-2 * time.Hour).In(moscowLocation), EndTime: time.Now().Add(2 * time.Hour).In(moscowLocation), Status: entity.StatusActive},
	}

	for _, d := range draws {
		_, err := db.Exec(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4)`, d.LotteryType, d.StartTime, d.EndTime, d.Status)
		require.NoError(t, err)
	}

	activeDraws, err := repo.GetActiveDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, activeDraws, 2)
}

func TestCancelDraw(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	startTime := time.Now().Add(-1 * time.Hour).In(moscowLocation)
	endTime := time.Now().Add(1 * time.Hour).In(moscowLocation)

	var drawID int32
	err := db.QueryRow(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4) RETURNING id`, "5 from 36", startTime, endTime, entity.StatusActive).Scan(&drawID)
	require.NoError(t, err)

	draw, err := repo.CancelDraw(context.Background(), drawID)
	assert.NoError(t, err)
	assert.Equal(t, entity.StatusCancelled, draw.Status)
	assert.Equal(t, drawID, draw.ID)
}

func TestActivateDraws(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	plannedDraws := []*entity.Draw{
		{LotteryType: "5 from 36", StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour), Status: entity.StatusPlanned},
		{LotteryType: "6 from 49", StartTime: now.Add(-2 * time.Hour), EndTime: now.Add(2 * time.Hour), Status: entity.StatusPlanned},
	}

	for _, d := range plannedDraws {
		_, err := db.Exec(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4)`, d.LotteryType, d.StartTime, d.EndTime, d.Status)
		require.NoError(t, err)
	}

	activatedDraws, err := repo.ActivateDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, activatedDraws, 2)
	for _, d := range activatedDraws {
		assert.Equal(t, entity.StatusActive, d.Status)
	}
}

func TestCompleteDraws(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	draws := []*entity.Draw{
		{LotteryType: "5 from 36", StartTime: now.Add(-3 * time.Hour), EndTime: now.Add(-1 * time.Hour), Status: entity.StatusActive},
		{LotteryType: "6 from 49", StartTime: now.Add(-4 * time.Hour), EndTime: now.Add(-2 * time.Hour), Status: entity.StatusActive},
	}

	for _, d := range draws {
		_, err := db.Exec(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4)`, d.LotteryType, d.StartTime, d.EndTime, d.Status)
		require.NoError(t, err)
	}

	completedDraws, err := repo.CompleteDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, completedDraws, 2)
	for _, d := range completedDraws {
		assert.Equal(t, entity.StatusCompleted, d.Status)
	}
}

func TestGetCompletedDraws(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	draws := []*entity.Draw{
		{LotteryType: "5 from 36", StartTime: time.Now().Add(-2 * time.Hour).In(moscowLocation), EndTime: time.Now().Add(-1 * time.Hour).In(moscowLocation), Status: entity.StatusCompleted},
		{LotteryType: "6 from 49", StartTime: time.Now().Add(-3 * time.Hour).In(moscowLocation), EndTime: time.Now().Add(-2 * time.Hour).In(moscowLocation), Status: entity.StatusCompleted},
	}

	for _, d := range draws {
		_, err := db.Exec(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4)`, d.LotteryType, d.StartTime, d.EndTime, d.Status)
		require.NoError(t, err)
	}

	result, err := repo.GetCompletedDraws(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetDrawResultByDrawID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewDrawRepository(db)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")

	startTime := time.Now().Add(-2 * time.Hour).In(moscowLocation)
	endTime := time.Now().Add(-1 * time.Hour).In(moscowLocation)

	var drawID int32
	err := db.QueryRow(`INSERT INTO draw.draws (lottery_type, start_time, end_time, status) VALUES ($1, $2, $3, $4) RETURNING id`, "5 from 36", startTime, endTime, entity.StatusCompleted).Scan(&drawID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO draw.draw_results (draw_id, winning_combination, result_time) VALUES ($1, $2, $3)`, drawID, "1,2,3,4,5", time.Now().In(moscowLocation))
	require.NoError(t, err)

	result, err := repo.GetDrawResult(context.Background(), drawID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, drawID, result.DrawID)
	assert.Equal(t, "1,2,3,4,5", result.WinningCombination)
}
