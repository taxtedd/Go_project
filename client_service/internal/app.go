package taxiApp

import (
	"client_service/internal/app/api/handlers"
	"client_service/internal/app/config"
	"client_service/internal/app/mongodb"
	"context"
	"go.uber.org/zap"
	"log"
)

type App struct {
	cfg     *config.Config
	handler *handlers.ClientHandler
	Logger  *zap.Logger
}

func (app *App) Run() {
	go func() {
		err := app.handler.Server.ListenAndServe()
		if err != nil {
			app.Logger.Error("Error starting server", zap.Error(err))
		}
	}()

	app.handler.Database.Close()
}

func (app *App) Stop(ctx context.Context) {
	if err := app.handler.Server.Shutdown(ctx); err != nil {
		app.Logger.Error("Error shutting down server", zap.Error(err))
	}
}

func NewApp(cfg *config.Config) *App {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	db, err := mongodb.NewDatabase(cfg.Database.URI, cfg.Database.Name) // Pass the database name from config
	if err != nil {
		logger.Error("Error initializing database", zap.Error(err))
	}

	app := &App{
		cfg:     cfg,
		handler: handlers.NewHandler(db, cfg, logger),
	}
	return app
}
