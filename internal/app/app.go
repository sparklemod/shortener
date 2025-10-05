package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	srv := httptransport.New(uc)

	return &App{server: srv, pg: pg}, nil
}

func (a *App) Run(ctx context.Context, cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: a.server.Router(),
	}

	go func() {
		log.Printf("starting HTTP server on %s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down HTTP server...")
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
		return err
	}

	log.Println("closing DB connection...")
	a.pg.Close()

	log.Println("server stopped")
	return nil
}
