package sqlxtransaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLX struct {
	db *sqlx.DB
}

func NewSQLX(db *sqlx.DB) SQLX {
	return SQLX{db: db}
}

func (s SQLX) Executor(ctx context.Context) Execer {
	txI := ctx.Value(txKey)
	if txI == nil {
		return s.db
	} else {
		return txI.(*sqlx.Tx)
	}
}

func (s SQLX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Executor(ctx).ExecContext(ctx, query, args...)
}

func (s SQLX) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.Executor(ctx).QueryxContext(ctx, query, args...)
}

func (s SQLX) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return s.Executor(ctx).QueryRowxContext(ctx, query, args...)
}

func (s SQLX) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.Executor(ctx).GetContext(ctx, dest, query, args...)
}

func (s SQLX) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.Executor(ctx).SelectContext(ctx, dest, query, args...)
}
