package server

import (
	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer(config *configs.Config, service *gofermart.Gofermart) *Server  {

	r := chi.NewRouter()

	return &Server{
		router: r,
	}
}

func (s* Server) Run() error{

	return nil
}

func (s *Server) Stop(err error) {

}


