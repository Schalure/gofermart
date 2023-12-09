package server

import (
	"net/http"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/server/handlers"
	"github.com/Schalure/gofermart/internal/server/midlewares"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer(config *configs.Config, handler *handlers.Handler, midleware *midlewares.Midleware) *Server  {

	r := chi.NewRouter()

	r.Post("/api/user/register", handler.UserRegistration)
	r.Post("/api/user/login", handler.UserAuthentication)

	return &Server{
		router: r,
	}
}

func (s* Server) Run(host string) error{

	return http.ListenAndServe(host, s.router)
}

func (s *Server) Stop(err error) {

}


