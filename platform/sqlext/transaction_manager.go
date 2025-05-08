package sqlext

import (
	"context"
	"database/sql"
	"github.com/MaxFando/lms/platform/sqlext/transaction"
	"time"

	"github.com/jmoiron/sqlx"
)

// TransactionManager предоставляет методы для управления транзакциями в базе данных, включая основные и вложенные транзакции.
// Использует sqlx для работы с базами данных и поддерживает обработку savepoint для внутренней транзакционности.
type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// RunTransaction выполняет переданную функцию в пределах транзакции. В случае ошибки инициирует откат изменений.
func (tm *TransactionManager) RunTransaction(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	tx, ok := transaction.GetTx(ctx)
	if ok {
		return tm.runInternalTransaction(ctx, tx, fn)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
		defer cancel()
	}

	o := new(sql.TxOptions)
	for _, opt := range opts {
		opt(o)
	}

	tx, err := tm.db.BeginTxx(ctx, o)
	if err != nil {
		return err
	}

	if err := fn(context.WithValue(ctx, transaction.TxKey, tx)); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

// runInternalTransaction выполняет операцию в рамках внутренней транзакции с использованием savepoint.
// Если переданная функция возвращает ошибку, происходит откат к savepoint.
// После успешного выполнения функции savepoint освобождается.
func (tm *TransactionManager) runInternalTransaction(ctx context.Context, tx *sqlx.Tx, fn transaction.AtomicFn) error {
	_subTx := newSubTx(tx)
	if err := _subTx.createSavepoint(ctx); err != nil {
		return err
	}

	if err := fn(ctx); err != nil {
		if rollbackErr := _subTx.rollbackSavepoint(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return _subTx.releaseSavepoint(ctx)
}
