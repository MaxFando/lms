package sqlext

import (
	"context"
	"database/sql"
	"fmt"
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
	tx, exists := transaction.GetTx(ctx)
	if exists {
		return tm.runInternalTransaction(ctx, tx, fn)
	}

	// Учет таймаутов в контексте
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) > 50*time.Second {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
		defer cancel()
	}

	// Настройка sql.TxOptions
	o := new(sql.TxOptions)
	for _, opt := range opts {
		opt(o)
	}

	// Начало транзакции
	tx, err := tm.db.BeginTxx(ctx, o)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Вызов переданной функции
	if err := fn(context.WithValue(ctx, transaction.TxKey, tx)); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error rolling back transaction: %w", rollbackErr)
		}
		return fmt.Errorf("error executing transaction: %w", err)
	}

	// Завершение транзакции
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
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
