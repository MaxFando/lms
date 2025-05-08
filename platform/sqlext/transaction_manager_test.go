package sqlext_test

import (
	"context"
	"github.com/MaxFando/lms/platform/sqlext"
	"github.com/MaxFando/lms/platform/sqlext/containers"
	"github.com/MaxFando/lms/platform/sqlext/transaction"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionManagerSuite struct {
	suite.Suite
	ctx         context.Context
	db          *sqlx.DB
	pgContainer *containers.PostgresContainer
}

func (suite *TransactionManagerSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := containers.CreatePostgresContainer(suite.ctx)
	suite.NoError(err)
	suite.pgContainer = pgContainer

	db, err := sqlx.ConnectContext(suite.ctx, "postgres", pgContainer.ConnectionString)
	suite.NoError(err)
	suite.db = db

	err = suite.db.PingContext(suite.ctx)
	suite.NoError(err)

	_, err = db.ExecContext(suite.ctx, "CREATE TABLE test_table (id serial PRIMARY KEY, name VARCHAR(50));")
	suite.NoError(err)
}

func (suite *TransactionManagerSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *TransactionManagerSuite) TearDownTest() {
	_, err := suite.db.ExecContext(suite.ctx, "TRUNCATE test_table;")
	suite.NoError(err)
}

func (suite *TransactionManagerSuite) TestSuccess() {
	tm := sqlext.NewTransactionManager(suite.db)

	err := tm.RunTransaction(suite.ctx, func(ctx context.Context) error {
		_, err := suite.db.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test');")
		return err
	})
	suite.NoError(err)

	var count int
	err = suite.db.GetContext(suite.ctx, &count, "SELECT COUNT(*) FROM test_table;")
	suite.NoError(err)
	suite.Equal(1, count)
}

func (suite *TransactionManagerSuite) TestRollback() {
	tm := sqlext.NewTransactionManager(suite.db)

	err := tm.RunTransaction(suite.ctx, func(ctx context.Context) error {
		executor := transaction.GetExecutor(ctx, suite.db)
		_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_rollback');")
		suite.NoError(err)

		return assert.AnError
	})
	suite.Error(err)

	var count int
	err = suite.db.GetContext(suite.ctx, &count, "SELECT COUNT(*) FROM test_table where name = 'test_rollback';")
	suite.NoError(err)
	suite.Equal(0, count)
}

func (suite *TransactionManagerSuite) TestSeveralNestedTransactions() {
	tm := sqlext.NewTransactionManager(suite.db)

	err := tm.RunTransaction(suite.ctx, func(ctx context.Context) error {
		executor := transaction.GetExecutor(ctx, suite.db)
		_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_1');")
		suite.NoError(err)

		return tm.RunTransaction(ctx, func(ctx context.Context) error {
			executor := transaction.GetExecutor(ctx, suite.db)
			_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_2');")
			suite.NoError(err)

			return tm.RunTransaction(ctx, func(ctx context.Context) error {
				executor := transaction.GetExecutor(ctx, suite.db)
				_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_3');")
				suite.NoError(err)

				return nil
			})
		})
	})
	suite.NoError(err)

	var count int
	err = suite.db.GetContext(suite.ctx, &count, "SELECT COUNT(*) FROM test_table;")
	suite.NoError(err)
	suite.Equal(3, count)
}

func (suite *TransactionManagerSuite) TestSeveralNestedTransactionsRollback() {
	tm := sqlext.NewTransactionManager(suite.db)

	err := tm.RunTransaction(suite.ctx, func(ctx context.Context) error {
		executor := transaction.GetExecutor(ctx, suite.db)
		_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_1');")
		suite.NoError(err)

		return tm.RunTransaction(ctx, func(ctx context.Context) error {
			executor := transaction.GetExecutor(ctx, suite.db)
			_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_2');")
			suite.NoError(err)

			return tm.RunTransaction(ctx, func(ctx context.Context) error {
				executor := transaction.GetExecutor(ctx, suite.db)
				_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_3');")
				suite.NoError(err)

				return assert.AnError
			})
		})
	})
	suite.Error(err)

	var count int
	err = suite.db.GetContext(suite.ctx, &count, "SELECT COUNT(*) FROM test_table;")
	suite.NoError(err)
	suite.Equal(0, count)
}

func (suite *TransactionManagerSuite) TestSeveralNestedTransactionsRollbackInner() {
	tm := sqlext.NewTransactionManager(suite.db)

	err := tm.RunTransaction(suite.ctx, func(ctx context.Context) error {
		executor := transaction.GetExecutor(ctx, suite.db)
		_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_inner_1');")
		suite.NoError(err)

		return tm.RunTransaction(ctx, func(ctx context.Context) error {
			executor := transaction.GetExecutor(ctx, suite.db)
			_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_inner_2');")
			suite.NoError(err)

			errInner := tm.RunTransaction(ctx, func(ctx context.Context) error {
				executor := transaction.GetExecutor(ctx, suite.db)
				_, err := executor.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test_nested_rollback_inner_3');")
				suite.NoError(err)

				return assert.AnError
			})

			suite.Error(errInner)
			return nil
		})
	})
	suite.NoError(err)

	var count int
	err = suite.db.GetContext(suite.ctx, &count, "SELECT COUNT(*) FROM test_table;")
	suite.NoError(err)
	suite.Equal(2, count)
}

func TestTransactionManagerSuite(t *testing.T) {
	suite.Run(t, new(TransactionManagerSuite))
}
