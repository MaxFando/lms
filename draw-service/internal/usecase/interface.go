package usecase

import (
	"context"

	"github.com/MaxFando/lms/draw-service/internal/entity"
)

// DrawRepository интерфейс для работы с репозиторием тиражей
type DrawRepository interface {

	// CreateDraw Создание нового тиража
	CreateDraw(ctx context.Context, draw *entity.Draw) (*entity.Draw, error)

	// GetActiveDraws Получение активных тиражей
	GetActiveDraws(ctx context.Context) ([]*entity.Draw, error)

	// CancelDraw Отмена тиража
	CancelDraw(ctx context.Context, id int32) (*entity.Draw, error)

	// ActivateDraws Активация тиражей
	ActivateDraws(ctx context.Context) ([]*entity.Draw, error)

	// CompleteDraws Завершение тиражей
	CompleteDraws(ctx context.Context) ([]*entity.Draw, error)

	// GetCompletedDraws Получение завершенных тиражей
	GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error)

	// GetDrawResult Получение инфы по завершенному тиражу
	GetDrawResult(ctx context.Context, drawID int32) (*entity.DrawResult, error)

	// BeginTransaction Старт транзакции
	BeginTransaction(ctx context.Context) (txContext context.Context, err error)

	// RollbackTransaction Отмена транзакции
	RollbackTransaction(txContext context.Context) (err error)

	// CommitTransaction Коммит транзакции
	CommitTransaction(txContext context.Context) (err error)
}

type DrawStatusQueue interface {
	PublishDraw(ctx context.Context, draw *entity.Draw, eventType entity.EventType) error
}
