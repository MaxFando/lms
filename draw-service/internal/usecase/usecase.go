package usecase

import (
	"context"
	"fmt"

	"github.com/MaxFando/lms/draw-service/internal/entity"
)

type DrawUseCase struct {
	drawRepo DrawRepository
}

func NewDrawUseCase(repo DrawRepository) *DrawUseCase {
	return &DrawUseCase{drawRepo: repo}
}

// CreateDraws - Создание нового тиража с типом лотереи и временем старта
func (uc *DrawUseCase) CreateDraws(ctx context.Context, draw entity.Draw) (*entity.Draw, error) {
	createdDraw, err := uc.drawRepo.CreateDraw(ctx, &draw)
	if err != nil {
		return nil, fmt.Errorf("create draw: %w", err)
	}

	return createdDraw, nil
}

// GetDrawsList - Получение списка активных тиражей
func (uc *DrawUseCase) GetDrawsList(ctx context.Context) ([]*entity.Draw, error) {
	draws, err := uc.drawRepo.GetActiveDraws(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active draws: %w", err)
	}
	return draws, nil
}

// CancelDraw - Отмена тиража (изменение статуса на CANCELLED)
func (uc *DrawUseCase) CancelDraw(ctx context.Context, id int64) error {
	err := uc.drawRepo.CancelDraw(ctx, id)
	if err != nil {
		return fmt.Errorf("cancel draw: %w", err)
	}

	return nil
}
