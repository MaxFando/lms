package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"

	"github.com/MaxFando/lms/platform/closer"
	"github.com/MaxFando/lms/platform/logger"
	"github.com/MaxFando/lms/platform/sqlext"
	"github.com/MaxFando/lms/platform/tracer"

	"github.com/MaxFando/lms/user-service/config"
	"github.com/MaxFando/lms/user-service/internal/controller"
	"github.com/MaxFando/lms/user-service/internal/jwt"
	pubsubPkg "github.com/MaxFando/lms/user-service/internal/pubsub"
	"github.com/MaxFando/lms/user-service/internal/repository"
	"github.com/MaxFando/lms/user-service/internal/server"
	"github.com/MaxFando/lms/user-service/internal/service"
)

type App struct {
	logger      logger.Logger
	config      *config.Config
	database    *sqlx.DB
	redisClient *redis.Client
	pubsub      pubsubPkg.PubSub
	srv         *server.Server
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

	// Redis + PubSub
	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     a.config.RedisAddr,
		Password: a.config.RedisPassword,
		DB:       a.config.RedisDB,
	})
	closer.Add(func() error { return a.redisClient.Close() })
	a.pubsub = pubsubPkg.NewRedisPubSub(a.redisClient)
	a.logger.Info(ctx, "Redis PubSub initialized")

	a.logger.Info(ctx, "Инициализация приложения завершена успешно")
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// репозиторий и сервис
	repo := repository.NewPostgresUserRepository(a.database)
	jwtSvc := jwt.NewJWTService(
		a.config.JWTSecret,
		a.config.AccessTokenTTL,
		a.config.RefreshTokenTTL,
	)
	userSvc := service.NewUserService(repo, jwtSvc, a.pubsub)

	userCtrl := controller.NewUserController(userSvc)
	srv := server.NewServer(a.logger, userCtrl)
	a.srv = srv
	go srv.Serve(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		a.logger.Info(ctx, "signal received, shutting down", "signal", sig.String())
	case err := <-srv.Notify():
		return fmt.Errorf("ошибка сервера: %w", err)
	}

	a.Shutdown(ctx)
	return nil
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
	cfg, err := tracer.NewConfig(
		a.config.TracerDSN,
		tracer.WithAppName(a.config.ServiceName),
		tracer.WithEnvironment(a.config.Env),
	)
	if err != nil {
		return fmt.Errorf("ошибка при создании конфигурации трейсинга: %w", err)
	}
	traceCloser, err := tracer.InitDefaultProvider(cfg)
	if err != nil {
		return fmt.Errorf("ошибка при инициализации провайдера трейсинга: %w", err)
	}
	closer.Add(func() error { return traceCloser(ctx) })
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
