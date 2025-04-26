package http

import (
	"net/http"

	"github.com/elvis-onobo/go-wallet-api/internal/http/handler"
	am "github.com/elvis-onobo/go-wallet-api/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go Wallet API"))
	})

	r.Post("/signup", handler.SignupHandler)
	r.Post("/login", handler.LoginHandler)
	r.Group(func(protected chi.Router) {
		protected.Use(am.JWTMiddleware)
		protected.Post("/wallet/fund", handler.FundWalletHandler)
		protected.Post("/wallet/withdraw", handler.WithdrawHandler)
	})

	return r
}
