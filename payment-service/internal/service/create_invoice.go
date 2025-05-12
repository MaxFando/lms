package service

import (
	"context"
	"time"

	"github.com/MaxFando/lms/payment-service/internal/entity"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateInvoice(ctx context.Context, userId int64, ticketId int64) (int64, decimal.Decimal, error) {
	ticket, err := s.ticket.BookTicket(ctx, userId, ticketId)
	if err != nil {
		return 0, decimal.Zero, err
	}

	registerTime := s.nowFunc()
	dueDate := registerTime.Add(15 * time.Minute)
	price := decimal.NewFromInt(s.cfg.TicketPrice)

	invoice := &entity.Invoice{
		Ticket:       ticket,
		OwnerID:      userId,
		Amount:       price,
		Status:       entity.InvoiceStatusPending,
		RegisterTime: registerTime,
		DueDate:      dueDate,
	}

	id, err := s.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		publishErr := s.publisher.PublishInvoice(ctx, invoice, entity.EventTypeInvoiceFailure)
		if publishErr != nil {
			s.log.Error(ctx, "failed to publish rollback task for ticked", ticketId, publishErr)
		}
		return 0, decimal.Zero, err
	}

	return id, price, nil
}

func (s *Service) CreateInvoiceForBookedTicket(ctx context.Context, userId int64, ticketId int64) (int64, decimal.Decimal, error) {
	registerTime := s.nowFunc()
	dueDate := registerTime.Add(15 * time.Minute)
	price := decimal.NewFromInt(s.cfg.TicketPrice)

	invoice := &entity.Invoice{
		Ticket: &entity.Ticket{
			ID: ticketId,
		},
		OwnerID:      userId,
		Amount:       price,
		Status:       entity.InvoiceStatusPending,
		RegisterTime: registerTime,
		DueDate:      dueDate,
	}

	id, err := s.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		return 0, decimal.Zero, err
	}

	return id, price, nil
}
