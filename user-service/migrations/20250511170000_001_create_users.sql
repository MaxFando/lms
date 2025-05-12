-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS consumer.users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    password      VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255),
    role          VARCHAR(32)  NOT NULL CHECK (role IN ('USER', 'ADMIN'))
);

CREATE INDEX IF NOT EXISTS users_name_idx ON consumer.users (name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS consumer.users;

-- +goose StatementEnd