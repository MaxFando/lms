package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/MaxFando/lms/draw-service/pkg/sqlxtransaction"
	"github.com/jmoiron/sqlx"
)

type DrawRepository struct {
	sqlxtransaction.SQLX
	sqlxtransaction.Transaction
}

func NewDrawRepository(db *sqlx.DB) *DrawRepository {
	return &DrawRepository{
		SQLX:        sqlxtransaction.NewSQLX(db),
		Transaction: sqlxtransaction.NewTransaction(db),
	}
}

// CreateDraw создает новый тираж с указанным типом лотереи и временем старта
func (r *DrawRepository) CreateDraw(ctx context.Context, draw *entity.Draw) (*entity.Draw, error) {
	query := `
		INSERT INTO draw.draws (lottery_type, start_time, end_time, status)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, lottery_type, start_time, end_time, status;
	`

	err := r.GetContext(ctx, draw, query, string(draw.LotteryType), draw.StartTime, draw.EndTime, entity.StatusPlanned)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return draw, nil
}

// GetActiveDraws возвращает список активных тиражей
func (r *DrawRepository) GetActiveDraws(ctx context.Context) ([]*entity.Draw, error) {
	query := `
		SELECT id, lottery_type, start_time, end_time, status 
		FROM draw.draws
		WHERE status = $1;
	`

	var draws []*entity.Draw
	err := r.SelectContext(ctx, &draws, query, entity.StatusActive)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return draws, nil
}

// CancelDraw изменяет статус тиража на CANCELLED
func (r *DrawRepository) CancelDraw(ctx context.Context, id int32) (*entity.Draw, error) {
	query := `
		UPDATE draw.draws
		SET status = 'CANCELLED'
		WHERE id = $1
		RETURNING id, lottery_type, start_time, end_time, status;
	`

	var draw entity.Draw
	err := r.GetContext(ctx, &draw, query, id)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &draw, nil
}

func (r *DrawRepository) ActivateDraws(ctx context.Context) ([]*entity.Draw, error) {
	query := `
		UPDATE draw.draws
		SET status = 'ACTIVE'
		WHERE status = 'PLANNED' AND start_time <= $1
		RETURNING id, lottery_type, start_time, end_time, status;
	`

	var updated []*entity.Draw
	err := r.SelectContext(ctx, &updated, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return updated, nil
}

// CompleteDueDraws sets status to COMPLETED for draws where end_time <= now and status = ACTIVE
func (r *DrawRepository) CompleteDraws(ctx context.Context) ([]*entity.Draw, error) {
	query := `
		UPDATE draw.draws
		SET status = 'COMPLETED'
		WHERE status = 'ACTIVE' AND end_time <= $1
		RETURNING id, lottery_type, start_time, end_time, status;
	`

	var updated []*entity.Draw
	err := r.SelectContext(ctx, &updated, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return updated, nil
}

// GetCompletedDraws возвращает все тиражи со статусом COMPLETED
func (r *DrawRepository) GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error) {
	query := `
		SELECT id, lottery_type, start_time, end_time, status
		FROM draw.draws
		WHERE status = 'COMPLETED';
	`

	var draws []*entity.Draw
	err := r.SelectContext(ctx, &draws, query)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return draws, nil
}

// GetDrawResultByDrawID возвращает результат тиража по draw_id
func (r *DrawRepository) GetDrawResult(ctx context.Context, drawID int32) (*entity.DrawResult, error) {
	query := `
		SELECT id, draw_id, winning_combination, result_time
		FROM draw.draw_results
		WHERE draw_id = $1;
	`

	var result entity.DrawResult
	err := r.GetContext(ctx, &result, query, drawID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &result, nil
}
