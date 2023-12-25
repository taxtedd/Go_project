package main

import (
	taxiApp "Go_project/internal/app"
	"Go_project/internal/app/config"
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
	flag.StringVar(&cfgPath, "c", ".config.yaml", "set path to config")
	flag.Parse()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(fmt.Errorf("fatal: init config #{err}"))
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
