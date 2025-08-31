package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HttpServer struct {
	validate *validator.Validate
	uc       UseCase
	router   *gin.Engine
}

func New(uc UseCase) *HttpServer {
	s := &HttpServer{
		validate: validator.New(),
		uc:       uc,
		router:   gin.Default(),
	}

	s.setupRoutes()
	return s
}

func (s *HttpServer) Router() *gin.Engine {
	return s.router
}

func (s *HttpServer) setupRoutes() {
	r := s.router
	api := r.Group("/api/v1")
	links := api.Group("/links")
	{
		links.POST("create", s.CreateLinkHandler)
	}
}
