package entity

import "time"

type Draw struct {
	ID          int32
	LotteryType string
	StartTime   time.Time
	EndTime     time.Time
	Status      string
}
