package service

import (
	"context"
	"errors"

	"github.com/MaxFando/lms/payment-service/internal/entity"
)

func (s *Service) Pay(ctx context.Context, userId int64, invoiceId int64, card *entity.Card) error {
	transactionID, err := s.payer.Pay(ctx, card)
	if err != nil {
		return err
	}

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	err = s.processPayment(tx, userId, invoiceId)
	if err != nil {
		s.processRollback(tx, transactionID, invoiceId)
		return err
	}

	return s.repo.CommitTransaction(tx)
}

func (s *Service) processPayment(ctx context.Context, userId int64, invoiceId int64) error {
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceId)
	if err != nil {
		return err
	}
	if invoice.OwnerID != userId {
		return errors.New("invoice is not owned by user")
	}

	err = s.repo.SetInvoiceStatus(ctx, invoiceId, entity.InvoiceStatusPaid)
	if err != nil {
		return err
	}

	_, err = s.repo.CreatePayment(ctx, invoiceId, entity.PaymentStatusPaid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) processRollback(ctx context.Context, transactionID int64, invoiceId int64) {
	rbErr := s.payer.Refund(ctx, transactionID)
	if rbErr != nil {
		s.log.Error(ctx, "failed to refund invoice for transaction", "transactionId", transactionID, "err", rbErr)
	}

	rbErr = s.repo.RollbackTransaction(ctx)
	if rbErr != nil {
		s.log.Error(ctx, "failed to rollback invoice", "invoiceID", invoiceId, "err", rbErr)
	}

	_, rbErr = s.repo.CreatePayment(ctx, invoiceId, entity.PaymentStatusRejected)
	if rbErr != nil {
		s.log.Error(ctx, "failed to create payment", "invoiceID", invoiceId, "err", rbErr)
	}
}
