//go:generate mockgen -source=$GOFILE -destination=./mock_${GOPACKAGE}_test.go -package=${GOPACKAGE}
package service

import (
	"context"
	"github.com/MaxFando/lms/draw-service/internal/core/draw/entity"
)

type DrawRepository interface {
	GetByID(ctx context.Context, drawID int32) (*entity.Draw, error)
}

type DrawService struct {
	repo DrawRepository
}

func NewDrawService(repo DrawRepository) *DrawService {
	return &DrawService{
		repo: repo,
	}
}

func (drawService *DrawService) GetByID(ctx context.Context, drawID int32) (*entity.Draw, error) {
	return drawService.repo.GetByID(ctx, drawID)
}
