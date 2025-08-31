package app

import (
	"context"
	"log"

	"shortener/config"
	adapterpg "shortener/internal/adapter/postgres"
	repopg "shortener/internal/repository/postgres"
	httptransport "shortener/internal/transport/http"
	"shortener/internal/usecase"
)

type App struct {
	server *httptransport.HttpServer
	//pg     *adapterpg.Postgres
}

func Build(cfg *config.Config) (*App, error) {
	ctx := context.Background()
	pg := adapterpg.New(cfg.LinksDB.URL)
	if err := pg.Connect(ctx); err != nil {
		return nil, err
	}
	repo := repopg.NewPostgres(pg)
	uc := usecase.New(repo)
	srv := httptransport.New(uc)
	//return &App{server: srv, pg: pg}, nil
	return &App{server: srv}, nil
}

func (a *App) Run(port string) error {
	log.Printf("starting HTTP server on %s", port)
	return a.server.Router().Run(port)
}

//func (a *App) Close(ctx context.Context) {
//	if a.pg != nil {
//		a.pg.Close(ctx)
//	}
//}
