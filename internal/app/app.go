package app

import (
	"context"
	"log"
	"shortener/config"
	adapterpg "shortener/internal/adapter/postgres"
	repopg "shortener/internal/repository"
	httptransport "shortener/internal/transport/http"
	"shortener/internal/usecase"
	"time"
)

type App struct {
	server *httptransport.HttpServer
	pg     *adapterpg.Postgres
}

func Build(ctx context.Context, cfg *config.Config) (*App, error) {
	pg := adapterpg.New(cfg.LinksDB.URL)
	if err := pg.Connect(ctx); err != nil {
		return nil, err
	}

	repo := repopg.NewPostgres(pg)
	uc := usecase.New(repo)
	srv := httptransport.New(uc, cfg)

	return &App{server: srv, pg: pg}, nil
}

func (a *App) Run() {
	a.server.Run()
}

func (a *App) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down HTTP server...")
	a.server.Stop(shutdownCtx)

	log.Println("closing DB connection...")
	a.pg.Close()

	log.Println("server stopped")
}
