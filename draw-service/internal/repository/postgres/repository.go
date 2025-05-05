package postgres

import (
	"context"
	"fmt"

	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/jmoiron/sqlx"
)

type DrawRepository struct {
	db *sqlx.DB
}

func NewDrawRepository(db *sqlx.DB) *DrawRepository {
	return &DrawRepository{
		db: db,
	}
}

// CreateDraw создает новый тираж с указанным типом лотереи и временем старта
func (r *DrawRepository) CreateDraw(ctx context.Context, draw *entity.Draw) (*entity.Draw, error) {
	query := `
		INSERT INTO draw_service.draw (lottery_type, start_time, end_time, status)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, lottery_type, start_time, end_time, status;
	`

	err := r.db.Get(draw, query, draw.LotteryType, draw.StartTime, draw.EndTime, "PLANNED")
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return draw, nil
}

// GetActiveDraws возвращает список активных тиражей
func (r *DrawRepository) GetActiveDraws(ctx context.Context) ([]*entity.Draw, error) {
	query := `
		SELECT id, lottery_type, start_time, end_time, status 
		FROM draw_service.draw
		WHERE status = $1;
	`

	var draws []*entity.Draw
	err := r.db.SelectContext(ctx, &draws, query, "ACTIVE")
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return draws, nil
}

// CancelDraw изменяет статус тиража на CANCELLED
func (r *DrawRepository) CancelDraw(ctx context.Context, id int64) error {
	query := `
		UPDATE draw_service.draw
		SET status = $1
		WHERE id = $2;
	`

	_, err := r.db.ExecContext(ctx, query, "CANCELLED", id)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
