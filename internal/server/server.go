package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer(handler *Handler, midleware *Middleware) *Server {

	r := chi.NewRouter()

	r.Post("/api/user/register", handler.UserRegistration)
	r.Post("/api/user/login", handler.UserAuthentication)

	r.Group(func(r chi.Router) {
		r.Use(midleware.WithAuthentication)
		r.Post("/api/user/orders", handler.LoadOrder)
		r.Get("/api/user/orders", handler.GetOrders)
		r.Get("/api/user/balance", handler.GetBalance)
		r.Post("/api/user/balance/withdraw", handler.WithdrawLoyaltyPoints)
		r.Get("/api/user/withdrawals", handler.GetOrdersWithdrawals)
	})

	return &Server{
		router: r,
	}
}

func (s *Server) Run(host string) error {

	return http.ListenAndServe(host, s.router)
}

func (s *Server) Stop(err error) {

}
