package postgres

import (
	"context"
	"github.com/MaxFando/lms/draw-service/internal/core/draw/entity"
	"github.com/jmoiron/sqlx"
)

type DrawRepositoryI interface {
	GetByID(ctx context.Context, drawID int32) (*entity.Draw, error)
}

type DrawRepository struct {
	db *sqlx.DB
}

func NewDrawRepository(db *sqlx.DB) *DrawRepository {
	return &DrawRepository{
		db: db,
	}
}

func (r *DrawRepository) GetByID(ctx context.Context, drawID int32) (*entity.Draw, error) {
	panic("implement me")
}
