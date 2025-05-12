package entity

type EventType string

const (
	EventTypeInvoiceOverdue EventType = "invoice_overdue"
	EventTypeInvoiceFailure EventType = "invoice_failure"
)
