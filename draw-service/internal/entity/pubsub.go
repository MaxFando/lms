package entity

type EventType string

const (
	EventTypeDrawActivated EventType = "draw_activated"
	EventTypeDrawCancelled EventType = "draw_cancelled"
	EventTypeDrawCompleted EventType = "draw_completed"
)
