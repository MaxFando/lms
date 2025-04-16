package sqlext

import (
	"context"
	"fmt"
	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5"
	pgxstd "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func OpenSqlxViaPgxConnPool(ctx context.Context, dsn string, opts ...ConnOption) (*sqlx.DB, error) {
	// Считываем конфигурацию подключения
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании конфигурации подключения: %w", err)
	}

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе DSN: %w", err)
	}
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	dsn = pgxstd.RegisterConnConfig(connConfig)

	attrs := append(otelsql.AttributesFromDSN(dsn), semconv.DBSystemPostgreSQL)
	options := []otelsql.Option{
		otelsql.WithAttributes(attrs...),
	}

	if cfg.tracerProvider != nil {
		options = append(options, otelsql.WithTracerProvider(cfg.tracerProvider))
	}

	wrappedDriverName, err := otelsql.Register("pgx", options...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при регистрации драйвера: %w", err)
	}

	// Открываем соединение
	dbConn, err := otelsql.Open(wrappedDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("подключение к базе данных: %w", err)
	}
	// ----------------------------

	// Настраиваем пул соединений
	dbConn.SetMaxIdleConns(cfg.maxIdleConns)
	dbConn.SetMaxOpenConns(cfg.maxOpenConns)

	dbConn.SetConnMaxLifetime(cfg.connLifeTime)
	dbConn.SetConnMaxIdleTime(cfg.connIdleTime)
	// ----------------------------

	// Проверяем соединение
	if err := dbConn.PingContext(ctx); err != nil {
		_ = dbConn.Close()
		return nil, fmt.Errorf("ошибка при пинге базы данных: %w", err)
	}
	// ----------------------------

	// Регистрируем метрики
	if err := otelsql.RegisterDBStatsMetrics(dbConn, otelsql.WithAttributes(attrs...)); err != nil {
		_ = dbConn.Close()
		return nil, fmt.Errorf("ошибка при регистрации метрик: %w", err)
	}
	// ----------------------------

	// Оборачиваем в sqlx
	return sqlx.NewDb(dbConn, wrappedDriverName), nil
}
