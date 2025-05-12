package app

import (
	"context"
	"fmt"
	"github.com/MaxFando/lms/platform/closer"
	"github.com/MaxFando/lms/platform/logger"
	"github.com/MaxFando/lms/platform/sqlext"
	"github.com/MaxFando/lms/platform/tracer"
	"github.com/MaxFando/lms/ticket-service/internal/repository/postgres"
	"github.com/MaxFando/lms/ticket-service/internal/service"
	"github.com/MaxFando/lms/ticket-service/internal/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"syscall"

	"github.com/MaxFando/lms/ticket-service/config"
	"github.com/MaxFando/lms/ticket-service/internal/server"
	v1 "github.com/MaxFando/lms/ticket-service/internal/server/service/v1"
)

type App struct {
	logger   logger.Logger
	config   *config.Config
	database *sqlx.DB
	srv      *server.Server
}

func New(cfg *config.Config) *App {
	l := logger.NewLogger()

	return &App{
		logger: l.With("app", "lms"),
		config: cfg,
	}
}

func (a *App) Logger() logger.Logger {
	return a.logger
}

func (a *App) Init(ctx context.Context) error {
	a.initCloser()

	if err := a.initTracer(ctx); err != nil {
		return fmt.Errorf("ошибка при инициализации трейсинга: %w", err)
	}

	if err := a.initDatabaseConnection(ctx); err != nil {
		return fmt.Errorf("ошибка при инициализации подключения к базе данных: %w", err)
	}

	a.logger.Info(ctx, "Инициализация приложения завершена успешно")

	return nil
}

func (a *App) Run(ctx context.Context) error {
	repo := postgres.NewTicketRepository(a.database)
	opt, err := redis.ParseURL(a.config.RedisDSN)
	if err != nil {
		return fmt.Errorf("redis parse url: %w", err)
	}

	rdb := redis.NewClient(opt)
	uc := usecase.NewTicketUsecase(repo)

	serviceServer := v1.NewServer(uc)
	srv := server.NewServer(a.logger, serviceServer)

	go func() {
		srv.Serve(ctx)
	}()

	a.srv = srv

	errChan := make(chan error, 1)

	drawHandler := service.NewDrawEventHandler(rdb, uc, a.config.RedisDrawChannel, 50)
	go func() {
		errChan <- drawHandler.Run(ctx)
	}()

	invoiceHandler := service.NewInvoiceEventHandler(rdb, uc, a.config.RedisInvoiceChannel)
	go func() {
		errChan <- invoiceHandler.Run(ctx)
	}()

	select {
	case s := <-srv.Notify():
		return fmt.Errorf("ошибка сервера: %w", s)
	case <-ctx.Done():
		return fmt.Errorf("ошибка контекста: %w", ctx.Err())
	case err := <-errChan:
		return fmt.Errorf("ошибка кроны: %w", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	defer func() {
		closer.CloseAll(ctx)
		closer.Wait()
	}()

	a.srv.Shutdown(ctx)
}

func (a *App) initCloser() {
	closer.New(syscall.SIGTERM, syscall.SIGINT)
}

func (a *App) initTracer(ctx context.Context) error {
	var tracingCfg tracer.Config
	tracingCfg, err := tracer.NewConfig(
		a.config.TracerDSN,
		tracer.WithAppName(a.config.ServiceName),
		tracer.WithEnvironment(a.config.Env),
	)

	if err != nil {
		return fmt.Errorf("ошибка при создании конфигурации трейсинга: %w", err)
	}

	var traceCloser tracer.ShutdownFn
	traceCloser, err = tracer.InitDefaultProvider(tracingCfg)
	if err != nil {
		return fmt.Errorf("ошибка при инициализации провайдера трейсинга: %w", err)
	}

	closer.Add(func() error {
		return traceCloser(ctx)
	})

	a.logger.Info(ctx, "Трейсинг инициализирован")

	return nil
}

func (a *App) initDatabaseConnection(ctx context.Context) error {
	db, err := sqlext.OpenSqlxViaPgxConnPool(ctx, a.config.DatabaseDSN, sqlext.WithTracerProvider(tracer.GetTraceProvider()))
	if err != nil {
		return fmt.Errorf("ошибка при создании подключения к базе данных: %w", err)
	}

	a.database = db
	closer.Add(func() error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("ошибка при закрытии подключения к базе данных: %w", err)
		}

		return nil
	})

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка при пинге базы данных: %w", err)
	}

	a.logger.Info(ctx, "Подключение к базе данных успешно установлено")

	return nil
}
