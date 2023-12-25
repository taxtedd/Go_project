package taxiApp

import (
	"Go_project/internal/app/api/handlers"
	"Go_project/internal/app/config"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	handler *handlers.OfferingHandler
	Logger  *zap.Logger
}

func (app *App) Run() {
	go func() {
		err := app.handler.Server.ListenAndServe()
		if err != nil {
			app.Logger.Error("Error starting server", zap.Error(err))
		}
	}()
}

func (app *App) Stop(ctx context.Context) {
	err := app.handler.Server.Shutdown(ctx)
	if err != nil {
		app.Logger.Error("Error shutting down server", zap.Error(err))
	}
}

func NewApp(cfg *config.Config) *App {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Logger init error.", err)
		return nil
	}
	app := &App{
		cfg:     cfg,
		handler: handlers.NewHandler(logger),
		Logger:  logger,
	}

	return app
}
