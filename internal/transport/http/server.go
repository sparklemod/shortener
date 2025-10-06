package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"shortener/config"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HttpServer struct {
	srv      *http.Server
	validate *validator.Validate
	uc       UseCase
}

func (h *HttpServer) Run() {
	go func() {
		log.Printf("starting HTTP server on %s", h.srv.Addr)
		if err := h.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (h *HttpServer) Stop(shutdownCtx context.Context) {
	if err := h.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}

func New(uc UseCase, cfg *config.Config) *HttpServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	s := &HttpServer{
		validate: validator.New(),
		uc:       uc,
		//router:   r,
		srv: &http.Server{
			Addr:    addr,
			Handler: gin.Default(),
		},
	}

	s.setupRoutes()
	return s
}

func (s *HttpServer) setupRoutes() {
	r := s.srv.Handler.(*gin.Engine)
	api := r.Group("/api/v1")
	links := api.Group("/links")
	{
		links.POST("create", s.CreateLinkHandler)
		links.GET("/:shorten-url", s.RedirectHandler)
	}
}
