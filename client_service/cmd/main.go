package main

import (
	taxiApp "client_service/internal"
	"client_service/internal/app/config"
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", ".configs.json", "set path to configs")
	flag.Parse()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(fmt.Errorf("fatal: init configs #{err}"))
	}

	app := taxiApp.NewApp(cfg)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	app.Stop(shutdownCtx)
}
