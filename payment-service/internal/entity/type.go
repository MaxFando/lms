package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "PENDING"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusOverdue InvoiceStatus = "OVERDUE"
)

type Invoice struct {
	ID           int64           `json:"id" db:"id"`
	Ticket       *Ticket         `json:"ticket_data" db:"ticket_data"`
	OwnerID      int64           `json:"owner_id" db:"owner_id"`
	Amount       decimal.Decimal `json:"amount" db:"amount"`
	Status       InvoiceStatus   `json:"status" db:"status"`
	RegisterTime time.Time       `json:"register_time" db:"register_time"`
	DueDate      time.Time       `json:"due_date" db:"due_date"`
}

type Ticket struct {
	ID int64 `json:"id" db:"id"`
}

func (t *Ticket) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for Ticket: %T", value)
	}
	return json.Unmarshal(bytes, t)
}

func (t *Ticket) Value() (driver.Value, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Ticket: %w", err)
	}
	return string(data), nil
}

type PaymentStatus string

const (
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusRejected PaymentStatus = "REJECTED"
)

type Payment struct {
	ID          int64         `json:"id" db:"id"`
	InvoiceID   int64         `json:"invoice_id" db:"invoice_id"`
	Status      PaymentStatus `json:"status" db:"status"`
	PaymentTime time.Time     `json:"payment_time" db:"payment_time"`
}

type Card struct {
	Number  string `json:"id"`
	ExpDate string `json:"invoice_id"`
	CVV     string `json:"status"`
}

func (c *Card) Validate() error {
	cardNumberRegex := regexp.MustCompile(`^\d{13,19}$`)
	if !cardNumberRegex.MatchString(c.Number) {
		return errors.New("invalid card number format")
	}

	cvvRegex := regexp.MustCompile(`^\d{3,4}$`)
	if !cvvRegex.MatchString(c.CVV) {
		return errors.New("invalid CVV format")
	}

	expDateRegex := regexp.MustCompile(`^(0[1-9]|1[0-2])\/\d{2}$`)
	if !expDateRegex.MatchString(c.ExpDate) {
		return errors.New("invalid expiration date format")
	}

	currentTime := time.Now()

	month := c.ExpDate[:2]
	year := c.ExpDate[3:]
	expiration, err := time.Parse("01/06", fmt.Sprintf("%s/%s", month, year))
	if err != nil {
		return errors.New("failed to parse expiration date")
	}

	endOfMonth := time.Date(expiration.Year(), expiration.Month()+1, 0, 23, 59, 59, 0, time.UTC)
	if currentTime.After(endOfMonth) {
		return errors.New("card has expired")
	}

	return nil
}
