-- +goose Up
-- +goose StatementBegin
CREATE TYPE ticket.ticket_status AS ENUM ('PENDING', 'WIN', 'LOSE');

CREATE TABLE IF NOT EXISTS ticket.tickets (
    ticket_id   SERIAL PRIMARY KEY,
    user_id     INT         NULL,
    draw_id     INT         NOT NULL,
    numbers     TEXT[]      NOT NULL,
    status      ticket.ticket_status NOT NULL DEFAULT 'PENDING',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_users   FOREIGN KEY(user_id) REFERENCES user.users(id),
    CONSTRAINT fk_draws   FOREIGN KEY(draw_id) REFERENCES draw.draws(id)
);

CREATE INDEX idx_tickets_user_draw ON ticket.tickets(user_id, draw_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ticket.tickets;
DROP TYPE IF EXISTS ticket.ticket_status;
-- +goose StatementEnd
