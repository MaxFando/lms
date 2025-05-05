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
	CancelDraw(ctx context.Context, id int64) error
}
