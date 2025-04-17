package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/MaxFando/lms/user-service/config"
	"github.com/MaxFando/lms/user-service/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()

	application := app.New(cfg)
	application.Logger().Info(ctx, "Инициализация приложения")

	if err := application.Init(ctx); err != nil {
		application.Logger().Error(ctx, "Завершение работы приложения с ошибкой", "error", err)

		return
	}

	if err := application.Run(ctx); err != nil {
		application.Logger().Error(ctx, "Завершение работы приложения с ошибкой", "error", err)
	}

	cancel()
	application.Logger().Info(ctx, "Завершение работы приложения")
	application.Shutdown(ctx)
	application.Logger().Info(ctx, "Завершение работы приложения завершено")
}
