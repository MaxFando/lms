-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS payment;

CREATE TYPE payment.invoice_status AS ENUM ('PENDING', 'PAID', 'OVERDUE');

CREATE TABLE IF NOT EXISTS payment.invoices (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    amount DECIMAL(12,2),
    ticket_data JSONB NOT NULL,
    status payment.invoice_status NOT NULL,
    register_time TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL
);

CREATE TYPE payment.payment_status AS ENUM ('PAID', 'REJECTED');

CREATE TABLE IF NOT EXISTS payment.payments (
    id SERIAL PRIMARY KEY,
    invoice_id INTEGER NOT NULL REFERENCES payment.invoices ON DELETE CASCADE,
    status payment.payment_status NOT NULL,
    payment_time TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment.invoices;
DROP TABLE IF EXISTS payment.payments;

DROP TYPE IF EXISTS payment.invoice_status;
DROP TYPE IF EXISTS payment.payment_status;

DROP SCHEMA IF EXISTS payment;
-- +goose StatementEnd

