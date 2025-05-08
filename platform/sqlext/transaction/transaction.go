package transaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// TxKeyType представляет собой пустую структуру, используемую для идентификации уникального ключа транзакции в контексте.
type TxKeyType struct{}

// TxKey используется как ключ для хранения и извлечения объекта транзакции из контекста.
var TxKey = &TxKeyType{}

// AtomicFn представляет собой функцию, выполняемую в пределах транзакции.
// Она принимает контекст и возвращает ошибку, если выполнение не удалось.
type AtomicFn func(ctx context.Context) error

type TxOption func(options *sql.TxOptions)

// WithTxIsolationLevel устанавливает определенный уровень изоляции для транзакции.
func WithTxIsolationLevel(level sql.IsolationLevel) TxOption {
	return func(options *sql.TxOptions) {
		options.Isolation = level
	}
}

type Executor interface {
	// ExecContext выполняет SQL-запрос без возврата строк, используя переданный контекст и параметры.
	// Возвращает объект sql.Result и ошибку, если таковая возникла.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// QueryContext выполняет SQL-запрос, принимает контекст исполнения, запрос и параметры, возвращает строки результата и ошибку.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// GetContext выполняет SQL-запрос для извлечения одного результата и заполняет его в переменную назначения dest.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// SelectContext выполняет запрос с использованием предоставленного контекста и заполняет результирующий объект dest.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// NamedExecContext выполняет запрос с использованием предоставленного контекста и именованных параметров.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// GetTx извлекает объект транзакции sqlx.Tx из контекста. Возвращает транзакцию и флаг её присутствия.
func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sqlx.Tx)
	return tx, ok
}

// GetExecutor возвращает Executor, используемый для выполнения запросов. Если в контексте задана транзакция, возвращает ее.
func GetExecutor(ctx context.Context, db *sqlx.DB) Executor {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}

	return db
}
