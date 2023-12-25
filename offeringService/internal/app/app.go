package taxiApp

import (
	"Go_project/offeringService/internal/app/api/handlers"
	"Go_project/offeringService/internal/app/config"
	"context"
	"fmt"
)

type App struct {
	cfg     *config.Config
	handler *handlers.OfferingHandler
}

func (app *App) Run() {
	go func() {
		err := app.handler.server.ListenAndServe()
		if err != nil {
			fmt.Println("Ошибка запуска сервера:", err)
		}
	}()
}

func (app *App) Stop(ctx context.Context) {
	fmt.Println(app.handler.server.Shutdown(ctx))
}

func NewApp(cfg *config.Config) *App {
	app := &App{
		cfg:     cfg,
		handler: handlers.NewHandler(),
	}

	return app
}
