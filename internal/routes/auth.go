package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/improver2108/jwt-login/internal/handler"
)

func RegisterAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Group(func(r chi.Router) {
		r.Post("/register", h.RegisterHandler)
		r.Post("/login", h.LoginHandler)
		r.Get("/logout", h.LogoutHandler)
		r.Get("/refresh-token", h.TokenRefreshHandler)
	})
}
