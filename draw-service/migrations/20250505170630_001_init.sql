-- +goose Up
-- +goose StatementBegin

CREATE TABLE draw.draw (
    id SERIAL PRIMARY KEY,
    lottery_type VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('PLANNED', 'ACTIVE', 'COMPLETED', 'CANCELLED'))
);

CREATE TABLE draw.draw_result (
    id SERIAL PRIMARY KEY,
    draw_id INTEGER NOT NULL REFERENCES draw.draw(id) ON DELETE CASCADE,
    winning_combination TEXT NOT NULL,
    result_time TIMESTAMPTZ NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS draw.draw_result;
DROP TABLE IF EXISTS draw.draw;
-- +goose StatementEnd

