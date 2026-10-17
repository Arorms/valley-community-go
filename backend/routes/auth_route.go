package routes

import (
	"valley-community-go/handlers"
	"valley-community-go/middlewares"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, authHandler *handlers.AuthHandler, auth *middlewares.TokenValidator) {
	r.Route("/api/auth", func(r chi.Router) {
		// GET /api/auth/me 需要携带 Authorization: Bearer <token>
		r.With(auth.Authenticate).Get("/me", authHandler.Me)
	})
}