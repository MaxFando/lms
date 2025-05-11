package usecase

import (
	"context"
	"fmt"

	"github.com/MaxFando/lms/platform/logger"

	"github.com/MaxFando/lms/draw-service/internal/entity"
)

type DrawUseCase struct {
	drawRepo  DrawRepository
	drawQueue DrawStatusQueue
	log       logger.Logger
}

func NewDrawUseCase(repo DrawRepository, queue DrawStatusQueue) *DrawUseCase {
	return &DrawUseCase{
		drawRepo:  repo,
		drawQueue: queue,
		log:       logger.NewLogger().With("app", "lms", "component", "draw-service", "layer", "usecase"),
	}
}

// CreateDraws - Создание нового тиража
func (uc *DrawUseCase) CreateDraws(ctx context.Context, draw entity.Draw) (*entity.Draw, error) {
	uc.log.Info(ctx, "creating draw", "lottery_type", draw.LotteryType, "start_time", draw.StartTime)

	createdDraw, err := uc.drawRepo.CreateDraw(ctx, &draw)
	if err != nil {
		uc.log.Error(ctx, "failed to create draw", "error", err)
		return nil, fmt.Errorf("create draw: %w", err)
	}

	uc.log.Info(ctx, "draw created successfully", "draw_id", createdDraw.ID)
	return createdDraw, nil
}

// GetDrawsList - Получение списка активных тиражей
func (uc *DrawUseCase) GetDrawsList(ctx context.Context) ([]*entity.Draw, error) {
	uc.log.Debug(ctx, "fetching active draws")

	draws, err := uc.drawRepo.GetActiveDraws(ctx)
	if err != nil {
		uc.log.Error(ctx, "failed to fetch active draws", "error", err)
		return nil, fmt.Errorf("get active draws: %w", err)
	}

	uc.log.Info(ctx, "active draws fetched", "count", len(draws))
	return draws, nil
}

// CancelDraw - Отмена тиража
func (uc *DrawUseCase) CancelDraw(ctx context.Context, id int32) error {
	log := uc.log.With("method", "CancelDraw", "draw_id", id)

	txCtx, err := uc.drawRepo.BeginTransaction(ctx)
	if err != nil {
		log.Error(ctx, "failed to begin transaction", "error", err)
		return fmt.Errorf("start transaction: %w", err)
	}
	defer uc.drawRepo.RollbackTransaction(txCtx)

	draw, err := uc.drawRepo.CancelDraw(txCtx, id)
	if err != nil {
		log.Error(ctx, "failed to cancel draw", "error", err, "draw_id", id)
		return fmt.Errorf("cancel draw: %w", err)
	}

	err = uc.drawQueue.PublishDraw(ctx, draw, entity.EventTypeDrawCancelled)
	if err != nil {
		log.Error(ctx, "failed to publish cancelled draw", "error", err, "draw_id", id)
		return fmt.Errorf("publish draw: %w", err)
	}

	err = uc.drawRepo.CommitTransaction(txCtx)
	if err != nil {
		log.Error(ctx, "failed to commit transaction", "error", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Info(ctx, "draw cancelled successfully", "draw_id", id)
	return nil
}

// MarkDrawsAsActive - Активация тиражей
func (uc *DrawUseCase) MarkDrawsAsActive(ctx context.Context) error {
	uc.log.Info(ctx, "updating draws statuses")

	txCtx, err := uc.drawRepo.BeginTransaction(ctx)
	if err != nil {
		uc.log.Error(ctx, "failed to begin transaction", "error", err)
		return fmt.Errorf("start transaction: %w", err)
	}
	defer uc.drawRepo.RollbackTransaction(txCtx)

	active, err := uc.drawRepo.ActivateDraws(txCtx)
	if err != nil {
		uc.log.Error(ctx, "failed to activate draws", "error", err)
		return fmt.Errorf("activate draws: %w", err)
	}

	for _, draw := range active {
		err := uc.drawQueue.PublishDraw(ctx, draw, entity.EventTypeDrawActivated)
		if err != nil {
			uc.log.Error(ctx, "failed to publish draw update", "draw_id", draw.ID, "error", err)

			return fmt.Errorf("publish draw: %w", err)
		}
	}

	err = uc.drawRepo.CommitTransaction(txCtx)
	if err != nil {
		uc.log.Error(ctx, "failed to commit transaction", "error", err)

		return fmt.Errorf("commit transaction: %w", err)
	}

	uc.log.Info(ctx, "draws updated successfully", "activated", len(active))

	return nil
}

// MarkDrawsAsCompleted - Завершение тиражей
func (uc *DrawUseCase) MarkDrawsAsCompleted(ctx context.Context) error {
	uc.log.Info(ctx, "updating draws statuses")

	txCtx, err := uc.drawRepo.BeginTransaction(ctx)
	if err != nil {
		uc.log.Error(ctx, "failed to begin transaction", "error", err)

		return fmt.Errorf("start transaction: %w", err)
	}
	defer uc.drawRepo.RollbackTransaction(txCtx)

	completed, err := uc.drawRepo.CompleteDraws(txCtx)
	if err != nil {
		uc.log.Error(ctx, "failed to complete draws", "error", err)
		return fmt.Errorf("complete draws: %w", err)
	}

	for _, draw := range completed {
		err := uc.drawQueue.PublishDraw(ctx, draw, entity.EventTypeDrawCompleted)
		if err != nil {
			uc.log.Error(ctx, "failed to publish draw update", "draw_id", draw.ID, "error", err)
			return fmt.Errorf("publish draw: %w", err)
		}
	}

	err = uc.drawRepo.CommitTransaction(txCtx)
	if err != nil {
		uc.log.Error(ctx, "failed to commit transaction", "error", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	uc.log.Info(ctx, "draws updated successfully", "completed", len(completed))
	return nil
}

// GetCompletedDraws - Получение завершенных тиражей
func (uc *DrawUseCase) GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error) {
	uc.log.Info(ctx, "fetching completed draws")

	draws, err := uc.drawRepo.GetCompletedDraws(ctx)
	if err != nil {
		uc.log.Error(ctx, "failed to get completed draws", "error", err)
		return nil, fmt.Errorf("get completed draws: %w", err)
	}

	uc.log.Info(ctx, "completed draws fetched", "count", len(draws))
	return draws, nil
}

// GetDrawResult - Получение результатов тиража
func (uc *DrawUseCase) GetDrawResult(ctx context.Context, id int32) (*entity.DrawResult, error) {
	uc.log.Info(ctx, "getting draw result", "draw_id", id)

	result, err := uc.drawRepo.GetDrawResult(ctx, id)
	if err != nil {
		uc.log.Error(ctx, "failed to get draw result", "draw_id", id, "error", err)
		return nil, fmt.Errorf("get draw result info: %w", err)
	}

	uc.log.Info(ctx, "draw result fetched", "draw_id", result.DrawID)
	return result, nil
}
