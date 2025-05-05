package entity

import "time"

// LotteryType - тип лотереи
type LotteryType string

const (
	LotteryType5from36 LotteryType = "5 from 36"
)

// DrawStatus - тип для статусов тиража
type DrawStatus string

const (
	StatusPlanned   DrawStatus = "PLANNED"
	StatusActive    DrawStatus = "ACTIVE"
	StatusCompleted DrawStatus = "COMPLETED"
	StatusCancelled DrawStatus = "CANCELLED"
)

// Draw - структура для описания тиража
type Draw struct {
	ID          int32       `json:"id" db:"id"`                     // Уникальный идентификатор тиража
	LotteryType LotteryType `json:"lottery_type" db:"lottery_type"` // Тип лотереи
	StartTime   time.Time   `json:"start_time" db:"start_time"`     // Дата и время начала тиража
	EndTime     time.Time   `json:"end_time" db:"end_time"`         // Дата и время завершения тиража
	Status      DrawStatus  `json:"status" db:"status"`             // Статус тиража (PLANNED, ACTIVE, COMPLETED, CANCELLED)
}

// DrawResult - структура для описания результата тиража
type DrawResult struct {
	ID                 int32     `json:"id"`                  // Уникальный идентификатор результата
	DrawID             int32     `json:"draw_id"`             // Ссылка на тираж
	WinningCombination string    `json:"winning_combination"` // Выигрышная комбинация чисел
	ResultTime         time.Time `json:"result_time"`         // Время определения результатов
}
