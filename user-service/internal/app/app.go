package app

import (
	"context"
	"fmt"
	"syscall"

	"github.com/MaxFando/lms/platform/closer"
	"github.com/MaxFando/lms/platform/logger"
	"github.com/MaxFando/lms/platform/sqlext"
	"github.com/MaxFando/lms/platform/tracer"
	"github.com/jmoiron/sqlx"

	"github.com/MaxFando/lms/user-service/config"
	"github.com/MaxFando/lms/user-service/internal/controller"
	"github.com/MaxFando/lms/user-service/internal/jwt"
	"github.com/MaxFando/lms/user-service/internal/repository"
	"github.com/MaxFando/lms/user-service/internal/server"
	"github.com/MaxFando/lms/user-service/internal/service"
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
	repo := repository.NewPostgresUserRepository(a.database)

	jwtSvc := jwt.NewJWTService(
		a.config.JWTSecret,
		a.config.AccessTokenTTL,
		a.config.RefreshTokenTTL,
	)

	userSvc := service.NewUserService(repo, jwtSvc)

	userCtrl := controller.NewUserController(userSvc)

	srv := server.NewServer(a.logger, userCtrl)
	srv.Serve(ctx)
	a.srv = srv

	select {
	case err := <-srv.Notify():
		return fmt.Errorf("ошибка сервера: %w", err)
	case <-ctx.Done():
		return fmt.Errorf("ошибка контекста: %w", ctx.Err())
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
	tracingCfg, err := tracer.NewConfig(
		a.config.TracerDSN,
		tracer.WithAppName(a.config.ServiceName),
		tracer.WithEnvironment(a.config.Env),
	)
	if err != nil {
		return fmt.Errorf("ошибка при создании конфигурации трейсинга: %w", err)
	}

	traceCloser, err := tracer.InitDefaultProvider(tracingCfg)
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
	db, err := sqlext.OpenSqlxViaPgxConnPool(
		ctx,
		a.config.DatabaseDSN,
		sqlext.WithTracerProvider(tracer.GetTraceProvider()),
	)
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