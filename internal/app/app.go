package app

import (
	"context"
	"log"
	"shortener/config"
	adaptermgo "shortener/internal/adapter/mongo"
	adapterpg "shortener/internal/adapter/postgres"
	"shortener/internal/repository/mongo"
	"shortener/internal/repository/postgres"
	httptransport "shortener/internal/transport/http"
	"shortener/internal/usecase"
	"time"
)

type App struct {
	server *httptransport.HttpServer
	dbConn DB
}

type DB interface {
	Close(ctx context.Context) error
}

func Build(ctx context.Context, cfg *config.Config) (*App, error) {
	var repo usecase.Repository
	var dbConn DB

	switch cfg.DBType {
	case config.Postgres:
		db := adapterpg.New(cfg.LinksDB.PostgresUrl)

		if err := db.RunMigrations(cfg.MigrationDir); err != nil {
			return nil, err
		}

		if err := db.Connect(ctx); err != nil {
			return nil, err
		}

		dbConn = db
		repo = postgres.NewPostgres(db)

	case config.Mongo:
		mg := adaptermgo.New(cfg.LinksDB.MongoUrl, cfg.DBName)
		if err := mg.Connect(ctx); err != nil {
			return nil, err
		}

		dbConn = mg
		repo = mongo.NewMongo(mg)
	}

	uc := usecase.New(repo)
	srv := httptransport.New(uc, cfg)

	return &App{server: srv, dbConn: dbConn}, nil
}

func (a *App) Run() {
	a.server.Run()
}

func (a *App) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down HTTP server...")
	if err := a.server.Stop(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}

	log.Println("closing DB connection...")
	if err := a.dbConn.Close(shutdownCtx); err != nil {
		log.Fatalf("error closing db: %v", err)
	}

	log.Println("server stopped")
}
