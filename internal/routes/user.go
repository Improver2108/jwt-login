package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/improver2108/jwt-login/internal/handler"
)

func RegisterUserRoutes(r chi.Router, h *handler.UserHandler) {
	r.Group(func(r chi.Router) {
	})
}
