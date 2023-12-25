package taxiApp

import (
	"client_service/internal/app/api/handlers"
	"client_service/internal/app/config"
	"client_service/internal/app/mongodb"
	"context"
	"fmt"
)

type App struct {
	cfg     *config.Config
	handler *handlers.ClientHandler
}

func (app *App) Run() {
	go func() {
		err := app.handler.Server.ListenAndServe()
		if err != nil {
			fmt.Println("Ошибка запуска сервера:", err)
		}
	}()

	app.handler.Database.Close()
}

func (app *App) Stop(ctx context.Context) {
	fmt.Println(app.handler.Server.Shutdown(ctx))
}

func NewApp(cfg *config.Config) *App {
	db, err := mongodb.NewDatabase(cfg.Database.URI, cfg.Database.Name) // Pass the database name from config
	if err != nil {
		fmt.Println("Ошибка инациализации бд:", err)
	}

	app := &App{
		cfg:     cfg,
		handler: handlers.NewHandler(db, cfg),
	}
	return app
}
