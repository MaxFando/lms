package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MaxFando/lms/ticket-service/internal/entity"
	"github.com/MaxFando/lms/ticket-service/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) repository.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) parseNumbersArray(arr string) []string {
	arr = strings.Trim(arr, "{}")
	if arr == "" {
		return []string{}
	}
	parts := strings.Split(arr, ",")
	for i := range parts {
		parts[i] = strings.Trim(parts[i], `"`)
	}
	return parts
}

func (r *TicketRepository) formatNumbersArray(nums []string) string {
	quoted := make([]string, len(nums))
	for i, s := range nums {
		quoted[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}"
}

func (r *TicketRepository) GetByID(ctx context.Context, id int32) (*entity.Ticket, error) {
	const query = `
        SELECT ticket_id, user_id, draw_id, numbers, status, created_at
        FROM ticket.tickets
        WHERE ticket_id = $1
    `
	var (
		t       entity.Ticket
		userID  sql.NullInt32
		numsArr string
		status  string
	)
	row := r.db.QueryRowxContext(ctx, query, id)
	if err := row.Scan(
		&t.ID,
		&userID,
		&t.DrawID,
		&numsArr,
		&status,
		&t.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ticket not found")
		}
		return nil, fmt.Errorf("scan ticket: %w", err)
	}
	if userID.Valid {
		uid := userID.Int32
		t.UserID = &uid
	}
	t.Numbers = r.parseNumbersArray(numsArr)
	t.Status = entity.Status(status)
	return &t, nil
}

func (r *TicketRepository) Create(ctx context.Context, t *entity.Ticket) (*entity.Ticket, error) {
	const query = `
        INSERT INTO ticket.tickets (user_id, draw_id, numbers, status)
        VALUES ($1, $2, $3::text[], $4)
        RETURNING ticket_id, created_at
    `
	numsLiteral := r.formatNumbersArray(t.Numbers)
	row := r.db.QueryRowxContext(ctx, query,
		t.UserID,
		t.DrawID,
		numsLiteral,
		string(t.Status),
	)
	if err := row.Scan(&t.ID, &t.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert ticket: %w", err)
	}
	return t, nil
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, id int32, status entity.Status) (*entity.Ticket, error) {
	const query = `
        UPDATE ticket.tickets
        SET status = $1
        WHERE ticket_id = $2
        RETURNING ticket_id, user_id, draw_id, numbers, status, created_at
    `
	var t entity.Ticket
	var numsArr string
	var st string
	row := r.db.QueryRowxContext(ctx, query, string(status), id)
	if err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.DrawID,
		&numsArr,
		&st,
		&t.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ticket not found")
		}
		return nil, fmt.Errorf("update status: %w", err)
	}
	t.Numbers = r.parseNumbersArray(numsArr)
	t.Status = entity.Status(st)
	return &t, nil
}

func (r *TicketRepository) ListByUser(ctx context.Context, userID int32) ([]*entity.TicketWithDraw, error) {
	const query = `
        SELECT
          t.ticket_id, t.user_id, t.draw_id, t.numbers, t.status, t.created_at,
          d.id, d.lottery_type, d.status, d.start_time, d.end_time AS draw_status
        FROM ticket.tickets t
        JOIN draw.draws d ON d.id = t.draw_id
        WHERE t.user_id = $1
        ORDER BY t.created_at DESC
    `
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query tickets: %w", err)
	}
	defer rows.Close()

	var out []*entity.TicketWithDraw
	for rows.Next() {
		var (
			t       entity.TicketWithDraw
			uID     sql.NullInt32
			numsArr string
			st      string
		)
		if err := rows.Scan(
			&t.ID, &uID, &t.DrawID, &numsArr, &st, &t.CreatedAt,
			&t.Draw.ID, &t.Draw.LotteryType, &t.Draw.Status, &t.Draw.StartTime, &t.Draw.EndTime,
		); err != nil {
			return nil, fmt.Errorf("scan ticket: %w", err)
		}
		if uID.Valid {
			u := uID.Int32
			t.UserID = &u
		}
		t.Numbers = r.parseNumbersArray(numsArr)
		t.Status = entity.Status(st)
		out = append(out, &t)
	}

	return out, rows.Err()
}

func (r *TicketRepository) IsDrawActive(ctx context.Context, drawID int32) (bool, error) {
	const query = `
        SELECT status
        FROM draw.draws
        WHERE id = $1
    `
	var st string
	row := r.db.QueryRowxContext(ctx, query, drawID)
	if err := row.Scan(&st); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("draw %d not found", drawID)
		}
		return false, fmt.Errorf("query draw status: %w", err)
	}
	return st == "ACTIVE", nil
}

func (r *TicketRepository) GetDrawLotteryType(ctx context.Context, drawID int32) (count int, maxNum int, err error) {
	const q = `SELECT lottery_type FROM draw.draws WHERE id = $1`
	var lt string
	if err = r.db.QueryRowxContext(ctx, q, drawID).Scan(&lt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, fmt.Errorf("draw %d not found", drawID)
		}
		return 0, 0, fmt.Errorf("scan draw config: %w", err)
	}

	parts := strings.Fields(lt)
	if len(parts) != 3 || parts[1] != "from" {
		return 0, 0, fmt.Errorf("invalid draw config format: %q", lt)
	}
	if count, err = strconv.Atoi(parts[0]); err != nil {
		return 0, 0, fmt.Errorf("parse count: %w", err)
	}
	if maxNum, err = strconv.Atoi(parts[2]); err != nil {
		return 0, 0, fmt.Errorf("parse max: %w", err)
	}
	return count, maxNum, nil
}

func (r *TicketRepository) BookTicket(ctx context.Context, ticketID, userID int32) (*entity.Ticket, error) {
	const q = `
        UPDATE ticket.tickets
        SET user_id = $1
        WHERE ticket_id = $2
          AND user_id IS NULL
        RETURNING ticket_id, user_id, draw_id, numbers, status, created_at
    `
	var (
		t       entity.Ticket
		uID     sql.NullInt32
		numsArr string
		st      string
	)
	row := r.db.QueryRowxContext(ctx, q, userID, ticketID)
	if err := row.Scan(&t.ID, &uID, &t.DrawID, &numsArr, &st, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("ticket %d is already booked or not found", ticketID)
		}
		return nil, fmt.Errorf("book ticket: %w", err)
	}
	if uID.Valid {
		u := uID.Int32
		t.UserID = &u
	}
	t.Numbers = r.parseNumbersArray(numsArr)
	t.Status = entity.Status(st)
	return &t, nil
}

func (r *TicketRepository) ClearBooking(ctx context.Context, ticketID int32) error {
	const q = `
        UPDATE ticket.tickets
        SET user_id = NULL, status = 'PENDING'
        WHERE ticket_id = $1
    `
	if _, err := r.db.ExecContext(ctx, q, ticketID); err != nil {
		return fmt.Errorf("clear booking: %w", err)
	}
	return nil
}

func (r *TicketRepository) ListFreeByActiveDraw(ctx context.Context) ([]*entity.Ticket, error) {
	const query = `
        SELECT
            t.ticket_id, t.user_id, t.draw_id, t.numbers, t.status, t.created_at
        FROM ticket.tickets t
        JOIN draw.draws d ON d.id = t.draw_id
        WHERE d.status = 'ACTIVE' AND t.user_id IS NULL
        ORDER BY t.created_at DESC
    `
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query free tickets: %w", err)
	}
	defer rows.Close()

	var tickets []*entity.Ticket
	for rows.Next() {
		var (
			t       entity.Ticket
			uID     sql.NullInt32
			numsArr string
			st      string
		)
		if err := rows.Scan(
			&t.ID, &uID, &t.DrawID, &numsArr, &st, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ticket: %w", err)
		}
		if uID.Valid {
			u := uID.Int32
			t.UserID = &u
		}
		t.Numbers = r.parseNumbersArray(numsArr)
		t.Status = entity.Status(st)
		tickets = append(tickets, &t)
	}
	return tickets, rows.Err()
}

func (r *TicketRepository) BulkUpdateStatus(ctx context.Context, ids []int32, status entity.Status) ([]*entity.Ticket, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)
	args[0] = string(status)
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}
	query := fmt.Sprintf(
		`UPDATE ticket.tickets
         SET status = $1
         WHERE ticket_id IN (%s)
         RETURNING ticket_id, user_id, draw_id, numbers, status, created_at`,
		strings.Join(placeholders, ","),
	)
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec bulk update: %w", err)
	}
	defer rows.Close()

	var tickets []*entity.Ticket
	for rows.Next() {
		var (
			t       entity.Ticket
			uID     sql.NullInt64
			numsArr string
			st      string
		)
		if err := rows.Scan(
			&t.ID, &uID, &t.DrawID, &numsArr, &st, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan updated ticket: %w", err)
		}
		if uID.Valid {
			u := int32(uID.Int64)
			t.UserID = &u
		}
		t.Numbers = r.parseNumbersArray(numsArr)
		t.Status = entity.Status(st)
		tickets = append(tickets, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate updated tickets: %w", err)
	}
	return tickets, nil
}
