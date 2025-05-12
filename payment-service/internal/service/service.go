package service

import (
	"context"
	"time"

	"github.com/MaxFando/lms/payment-service/config"
	"github.com/MaxFando/lms/payment-service/internal/entity"
	"github.com/MaxFando/lms/platform/logger"
)

type repo interface {
	CreateInvoice(ctx context.Context, invoice *entity.Invoice) (int64, error)
	GetInvoiceByID(ctx context.Context, id int64) (*entity.Invoice, error)
	GetPendingInvoices(ctx context.Context) ([]*entity.Invoice, error)
	SetInvoiceStatus(ctx context.Context, id int64, status entity.InvoiceStatus) error
	CreatePayment(ctx context.Context, invoiceID int64, status entity.PaymentStatus) (int64, error)
	BeginTransaction(ctx context.Context) (txContext context.Context, err error)
	RollbackTransaction(txContext context.Context) (err error)
	CommitTransaction(txContext context.Context) (err error)
}

type publisher interface {
	PublishInvoice(ctx context.Context, invoice *entity.Invoice, eventType entity.EventType) error
}

type ticketService interface {
	BookTicket(ctx context.Context, userId int64, ticketI int64) (*entity.Ticket, error)
}

type payer interface {
	Pay(ctx context.Context, card *entity.Card) (int64, error)
	Refund(ctx context.Context, transactionID int64) error
}

type Service struct {
	ticket    ticketService
	payer     payer
	repo      repo
	publisher publisher
	log       logger.Logger
	cfg       *config.Config

	nowFunc func() time.Time
}

func New(
	ticket ticketService,
	payer payer,
	repo repo,
	publisher publisher,
	cfg *config.Config,
) *Service {
	return &Service{
		ticket:    ticket,
		payer:     payer,
		repo:      repo,
		publisher: publisher,
		log:       logger.NewLogger().With("app", "lms", "component", "payment-service", "layer", "usecase"),
		cfg:       cfg,

		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}
