package sqlxtransaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) Transaction {
	return Transaction{db: db}
}

func (t Transaction) BeginTransaction(ctx context.Context) (txContext context.Context, err error) {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("beign: %w", err)
		return
	}

	return context.WithValue(ctx, txKey, tx), nil
}

func (t Transaction) RollbackTransaction(txContext context.Context) (err error) {
	txI := txContext.Value(txKey)
	if txI == nil {
		panic("no transaction in context")
	}
	err = txI.(*sqlx.Tx).Rollback()
	if err != nil {
		err = fmt.Errorf("rollback: %w", err)
	}
	return
}

func (t Transaction) CommitTransaction(txContext context.Context) (err error) {
	txI := txContext.Value(txKey)
	if txI == nil {
		panic("no transaction in context")
	}
	err = txI.(*sqlx.Tx).Commit()
	if err != nil {
		err = fmt.Errorf("commit: %w", err)
	}
	return
}
