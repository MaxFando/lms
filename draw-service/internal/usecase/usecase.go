package usecase

import (
	"context"
	"fmt"

	"github.com/MaxFando/lms/draw-service/internal/entity"
)

type DrawUseCase struct {
	drawRepo  DrawRepository
	drawQueue DrawStatusQueue
}

func NewDrawUseCase(repo DrawRepository, queue DrawStatusQueue) *DrawUseCase {
	return &DrawUseCase{drawRepo: repo, drawQueue: queue}
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
func (uc *DrawUseCase) CancelDraw(ctx context.Context, id int32) error {
	txCtx, err := uc.drawRepo.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer uc.drawRepo.RollbackTransaction(txCtx)

	draw, err := uc.drawRepo.CancelDraw(txCtx, id)
	if err != nil {
		return fmt.Errorf("cancel draw: %w", err)
	}

	err = uc.drawQueue.PublishDraw(ctx, draw)
	if err != nil {
		return fmt.Errorf("publish draw: %w", err)
	}

	err = uc.drawRepo.CommitTransaction(txCtx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// UpdateDraws - обновление статусов тиражей
func (uc *DrawUseCase) UpdateDraws(ctx context.Context) error {
	txCtx, err := uc.drawRepo.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer uc.drawRepo.RollbackTransaction(txCtx)

	draws, err := uc.drawRepo.ActivateDraws(txCtx)
	if err != nil {
		return fmt.Errorf("activate draws: %w", err)
	}

	drawsCompleted, err := uc.drawRepo.CompleteDraws(txCtx)
	if err != nil {
		return fmt.Errorf("complete draws: %w", err)
	}

	draws = append(draws, drawsCompleted...)

	for _, v := range draws {
		err := uc.drawQueue.PublishDraw(ctx, v)
		if err != nil {
			return fmt.Errorf("publish draw: %w", err)
		}
	}

	err = uc.drawRepo.CommitTransaction(txCtx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetCompletedDraws - получение завершенных тиражей
func (uc *DrawUseCase) GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error) {
	draws, err := uc.drawRepo.GetCompletedDraws(ctx)
	if err != nil {
		return nil, fmt.Errorf("get completed draws: %w", err)
	}
	return draws, nil
}

// GetDrawResult - получение информации по завершенному тиражу
func (uc *DrawUseCase) GetDrawResult(ctx context.Context, id int32) (*entity.DrawResult, error) {
	draw, err := uc.drawRepo.GetDrawResult(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get draw result info: %w", err)
	}
	return draw, nil
}
