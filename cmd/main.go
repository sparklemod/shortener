package main

import (
	"log"
	"shortener/config"
	"shortener/internal/app"
)

func main() {
	cfg := config.NewConfig()
	application, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}
	if err := application.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
