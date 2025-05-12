package repository

import (
	"context"
	"github.com/MaxFando/lms/ticket-service/internal/entity"
)

type TicketRepository interface {
	GetByID(ctx context.Context, id int32) (*entity.Ticket, error)
	Create(ctx context.Context, t *entity.Ticket) (*entity.Ticket, error)
	UpdateStatus(ctx context.Context, id int32, status entity.Status) (*entity.Ticket, error)
	ListByUser(ctx context.Context, userID int32) ([]*entity.TicketWithDraw, error)
	IsDrawActive(ctx context.Context, drawID int32) (bool, error)
	GetDrawLotteryType(ctx context.Context, drawID int32) (count int, maxNum int, err error)
	BookTicket(ctx context.Context, ticketID, userID int32) (*entity.Ticket, error)
	ClearBooking(ctx context.Context, ticketID int32) error
	ListFreeByActiveDraw(ctx context.Context) ([]*entity.Ticket, error)
	BulkUpdateStatus(ctx context.Context, ids []int32, status entity.Status) ([]*entity.Ticket, error)
}
