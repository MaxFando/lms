package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MaxFando/lms/payment-service/internal/entity"
	"github.com/MaxFando/lms/payment-service/pkg/sqlxtransaction"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
	sqlxtransaction.SQLX
	sqlxtransaction.Transaction
}

func New(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{
		SQLX:        sqlxtransaction.NewSQLX(db),
		Transaction: sqlxtransaction.NewTransaction(db),
	}
}

func (r *PaymentRepository) CreateInvoice(ctx context.Context, invoice *entity.Invoice) (int64, error) {
	query := `
		INSERT INTO payment.invoices (owner_id, amount, ticket_data, status, register_time, due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := r.GetContext(ctx, &id, query,
		invoice.OwnerID,
		invoice.Amount,
		invoice.Ticket,
		invoice.Status,
		invoice.RegisterTime,
		invoice.DueDate,
	)
	if err != nil {
		return 0, fmt.Errorf("create payment: %w", err)
	}

	return id, nil
}

func (r *PaymentRepository) GetPendingInvoices(ctx context.Context) ([]*entity.Invoice, error) {
	query := `
		SELECT id, owner_id, amount, ticket_data, status, register_time, due_date
		FROM payment.invoices
		WHERE status = 'PENDING'
	`
	var invoices []*entity.Invoice
	if err := r.SelectContext(ctx, &invoices, query); err != nil {
		return nil, fmt.Errorf("get pending invoices: %w", err)
	}

	return invoices, nil
}

func (r *PaymentRepository) GetInvoiceByID(ctx context.Context, id int64) (*entity.Invoice, error) {
	query := `
		SELECT id, owner_id, amount, ticket_data, status, register_time, due_date
		FROM payment.invoices
		WHERE id = &1
	`
	var invoices entity.Invoice
	if err := r.SelectContext(ctx, &invoices, query, id); err != nil {
		return nil, fmt.Errorf("get pending invoices: %w", err)
	}

	return &invoices, nil
}

func (r *PaymentRepository) SetInvoiceStatus(ctx context.Context, id int64, status entity.InvoiceStatus) error {
	query := `
		UPDATE payment.invoices
		SET status = $1
		WHERE id = $2
	`
	if _, err := r.ExecContext(ctx, query, status, id); err != nil {
		return fmt.Errorf("set invoice status: %w", err)
	}

	return nil
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, invoiceID int64, status entity.PaymentStatus) (int64, error) {
	query := `
		INSERT INTO payment.payments (invoice_id, status, payment_time)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64
	err := r.GetContext(ctx, &id, query,
		invoiceID,
		status,
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("create payment: %w", err)
	}
	return id, nil
}

func (r *PaymentRepository) GetInvoicePayments(ctx context.Context) ([]*entity.Payment, error) {
	query := `
		SELECT id, invoice_id, status, payment_time
		FROM payment.payments
		WHERE id IN (SELECT invoice_id FROM payment.payments)
	`
	var payments []*entity.Payment
	if err := r.SelectContext(ctx, &payments, query); err != nil {
		return nil, fmt.Errorf("get invoice payments: %w", err)
	}

	return payments, nil
}

func (r *PaymentRepository) SetPaymentStatus(ctx context.Context, id int64, status entity.PaymentStatus) error {
	query := `
		UPDATE payment.payments
		SET status = $1
		WHERE id = $2
	`
	if _, err := r.ExecContext(ctx, query, status, id); err != nil {
		return fmt.Errorf("set payment status: %w", err)
	}

	return nil
}
