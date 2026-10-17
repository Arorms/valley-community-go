package routes

import (
	"valley-community-go/handlers"
	"valley-community-go/middlewares"

	"github.com/go-chi/chi/v5"
)

func RegisterPostRoutes(r chi.Router, postHandler *handlers.PostHandler, auth *middlewares.TokenValidator) {
	r.Route("/api/post", func(r chi.Router) {
		r.Get("/", postHandler.GetByID)  // GET  /api/post?id=1
		r.Get("/page", postHandler.List) // GET  /api/post/page?offset=0&limit=10

		// POST /api/post 需要携带 Authorization: Bearer <token>
		r.With(auth.Authenticate).Post("/", postHandler.CreatePost)
	})
}