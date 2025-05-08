package sqlext

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// subTx представляет объект для работы с savepoint внутри транзакции SQL.
// Используется для управления savepoint и операций над ними.
type subTx struct {
	tx            *sqlx.Tx
	savepointName string
}

func newSubTx(tx *sqlx.Tx) *subTx {
	return &subTx{tx: tx, savepointName: generateSavepointName()}
}

// createSavepoint создает savepoint в текущей транзакции, используя заданное имя savepoint.
func (s *subTx) createSavepoint(ctx context.Context) error {
	_, err := s.tx.ExecContext(ctx, "SAVEPOINT "+s.savepointName)
	return err
}

// releaseSavepoint освобождает ранее созданный savepoint в рамках текущей транзакции.
// Возвращает ошибку, если выполнение команды RELEASE SAVEPOINT не удалось.
func (s *subTx) releaseSavepoint(ctx context.Context) error {
	_, err := s.tx.ExecContext(ctx, "RELEASE SAVEPOINT "+s.savepointName)
	return err
}

// rollbackSavepoint откатывает состояние транзакции к указанному savepoint. Возвращает ошибку в случае неудачи.
func (s *subTx) rollbackSavepoint(ctx context.Context) error {
	_, err := s.tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+s.savepointName)
	return err
}

// generateSavepointName генерирует уникальное имя для сохранённой точки на основе текущего времени в нано-секундах.
func generateSavepointName() string {
	return fmt.Sprintf("tx_sp_%d", time.Now().UnixNano())
}
