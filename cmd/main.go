package main

import (
	"context"
	"log"
	"os/signal"
	"shortener/config"
	"shortener/internal/app"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()
	application, err := app.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}
	application.Run()

	<-ctx.Done()
	application.Stop()
}
