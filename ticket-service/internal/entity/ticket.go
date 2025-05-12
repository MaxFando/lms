package entity

import "time"

type Status string

const (
	StatusPending Status = "PENDING"
	StatusWin     Status = "WIN"
	StatusLose    Status = "LOSE"
)

type Ticket struct {
	ID        int32
	UserID    *int32
	DrawID    int32
	Numbers   []string
	Status    Status
	CreatedAt time.Time
}

type TicketWithDraw struct {
	ID        int32
	UserID    *int32
	DrawID    int32
	Numbers   []string
	Status    Status
	CreatedAt time.Time
	Draw      Draw
}
