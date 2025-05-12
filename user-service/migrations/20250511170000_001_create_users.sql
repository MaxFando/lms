-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255)    NOT NULL UNIQUE,
    password       VARCHAR(255)    NOT NULL,
    refresh_token  VARCHAR(255),
    role           VARCHAR(32)     NOT NULL CHECK (role IN ('USER', 'ADMIN'))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd