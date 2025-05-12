package service

import (
	"context"
	"fmt"
	"github.com/MaxFando/lms/payment-service/internal/entity"
)

func (s *Service) ProcessInvoices(ctx context.Context) error {
	pendingInvoices, err := s.repo.GetPendingInvoices(ctx)
	if err != nil {
		return err
	}

	now := s.nowFunc()
	for _, invoice := range pendingInvoices {
		if invoice.DueDate.Before(now) {
			continue
		}

		err = s.processInvoiceTx(ctx, invoice)
		if err != nil {
			s.log.Error(ctx, err.Error())
		}
	}

	return nil
}

func (s *Service) processInvoiceTx(ctx context.Context, invoice *entity.Invoice) error {
	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	err = s.processInvoice(tx, invoice)
	if err != nil {
		txErr := s.repo.RollbackTransaction(tx)
		if txErr != nil {
			s.log.Error(ctx, txErr.Error())
		}
		return err
	}

	return s.repo.CommitTransaction(tx)
}

func (s *Service) processInvoice(ctx context.Context, invoice *entity.Invoice) error {
	err := s.repo.SetInvoiceStatus(ctx, invoice.ID, entity.InvoiceStatusOverdue)
	if err != nil {
		return fmt.Errorf("failded to set invoice status to overdue: %w", err)
	}

	err = s.publisher.PublishInvoice(ctx, invoice, entity.EventTypeInvoiceOverdue)
	if err != nil {
		return fmt.Errorf("failded to publish rollback task for invoice: %w", err)
	}

	return nil
}
