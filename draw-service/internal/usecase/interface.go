package usecase

import (
	"context"

	"github.com/MaxFando/lms/draw-service/internal/entity"
)

// DrawRepository интерфейс для работы с репозиторием тиражей
type DrawRepository interface {

	// Создание нового тиража
	CreateDraw(ctx context.Context, draw *entity.Draw) (*entity.Draw, error)

	// Получение активных тиражей
	GetActiveDraws(ctx context.Context) ([]*entity.Draw, error)

	// Отмена тиража
	CancelDraw(ctx context.Context, id int32) (*entity.Draw, error)

	// Активация тиражей
	ActivateDraws(ctx context.Context) ([]*entity.Draw, error)

	// Завершение тиражей
	CompleteDraws(ctx context.Context) ([]*entity.Draw, error)

	// Получение завершенных тиражей
	GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error)

	// Получение инфы по завершенному тиражу
	GetDrawResult(ctx context.Context, drawID int32) (*entity.DrawResult, error)

	// Старт транзакции
	BeginTransaction(ctx context.Context) (txContext context.Context, err error)

	// Отмена транзакции
	RollbackTransaction(txContext context.Context) (err error)

	// Коммит транзакции
	CommitTransaction(txContext context.Context) (err error)
}

type DrawStatusQueue interface {
	PublishDraw(ctx context.Context, draw *entity.Draw) error
}
