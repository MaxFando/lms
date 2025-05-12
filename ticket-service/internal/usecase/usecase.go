package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaxFando/lms/ticket-service/pkg/lottery"
	"strconv"
	"time"

	"github.com/MaxFando/lms/ticket-service/internal/entity"
	"github.com/MaxFando/lms/ticket-service/internal/repository"
)

var (
	ErrDrawNotActive  = errors.New("draw not active")
	ErrInvalidNumbers = errors.New("invalid ticket numbers")
)

type TicketUsecase struct {
	repo repository.TicketRepository
}

func NewTicketUsecase(repo repository.TicketRepository) *TicketUsecase {
	return &TicketUsecase{
		repo: repo,
	}
}

func (u *TicketUsecase) GetTicket(ctx context.Context, id int32) (*entity.Ticket, error) {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get ticket: %w", err)
	}
	return t, nil
}

func (u *TicketUsecase) CreateTicket(ctx context.Context, userID, drawID int32, numbers []string) (*entity.Ticket, error) {
	active, err := u.repo.IsDrawActive(ctx, drawID)
	if err != nil {
		return nil, fmt.Errorf("check draw active: %w", err)
	}
	if !active {
		return nil, ErrDrawNotActive
	}

	count, maxNum, err := u.repo.GetDrawLotteryType(ctx, drawID)
	if err != nil {
		return nil, fmt.Errorf("get draw lottery type: %w", err)
	}

	if len(numbers) != count {
		return nil, ErrInvalidNumbers
	}
	seen := make(map[string]struct{}, count)
	for _, s := range numbers {
		if _, dup := seen[s]; dup {
			return nil, ErrInvalidNumbers
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > maxNum {
			return nil, ErrInvalidNumbers
		}
		seen[s] = struct{}{}
	}

	ticket := &entity.Ticket{
		UserID:    &userID,
		DrawID:    drawID,
		Numbers:   numbers,
		Status:    entity.StatusPending,
		CreatedAt: time.Now(),
	}
	saved, err := u.repo.Create(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}

	// todo invoice

	return saved, nil
}

func (u *TicketUsecase) ReserveTicket(ctx context.Context, id int32) (*entity.Ticket, error) {
	t, err := u.repo.UpdateStatus(ctx, id, entity.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("reserve update ticket: %w", err)
	}
	return t, nil
}

func (u *TicketUsecase) ListUserTickets(ctx context.Context, userID int32) ([]*entity.TicketWithDraw, error) {
	tickets, err := u.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list repo: %w", err)
	}
	return tickets, nil
}

func (u *TicketUsecase) BookTicket(ctx context.Context, userID, ticketID int32) (*entity.Ticket, error) {
	booked, err := u.repo.BookTicket(ctx, ticketID, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase book ticket: %w", err)
	}
	return booked, nil
}

func (u *TicketUsecase) ReleaseBooking(ctx context.Context, ticketID int32) error {
	if err := u.repo.ClearBooking(ctx, ticketID); err != nil {
		return fmt.Errorf("usecase release booking: %w", err)
	}
	return nil
}

func (u *TicketUsecase) GenerateTickets(ctx context.Context, drawID int32, count int) error {
	needed, maxNum, err := u.repo.GetDrawLotteryType(ctx, drawID)
	if err != nil {
		return fmt.Errorf("get draw config: %w", err)
	}

	now := time.Now()
	for i := 0; i < count; i++ {
		nums, err := lottery.GenerateTicketNumbers(needed, maxNum)
		if err != nil {
			return fmt.Errorf("generate numbers for ticket %d: %w", i, err)
		}
		t := &entity.Ticket{
			UserID:    nil,
			DrawID:    drawID,
			Numbers:   nums,
			Status:    entity.StatusPending,
			CreatedAt: now,
		}
		if _, err := u.repo.Create(ctx, t); err != nil {
			return fmt.Errorf("create ticket %d: %w", i, err)
		}
	}
	return nil
}

func (u *TicketUsecase) ListAvailableTickets(ctx context.Context) ([]*entity.Ticket, error) {
	return u.repo.ListFreeByActiveDraw(ctx)
}

func (u *TicketUsecase) SetWinningTickets(ctx context.Context, ids []int32) ([]*entity.Ticket, error) {
	if len(ids) == 0 {
		return nil, ErrInvalidNumbers
	}
	return u.repo.BulkUpdateStatus(ctx, ids, entity.StatusWin)
}

func (u *TicketUsecase) CheckResult(ctx context.Context, ticketID int32) (string, error) {
	t, err := u.repo.GetByID(ctx, ticketID)
	if err != nil {
		return "", err
	}
	return string(t.Status), nil
}
