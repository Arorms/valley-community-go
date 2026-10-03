package routes

import (
	"valley-community-go/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterPostRoutes(r chi.Router, postHandler *handlers.PostHandler) {
	r.Route("/api/post", func(r chi.Router) {
		r.Post("/", postHandler.CreatePost) // POST /api/posts
		r.Get("/", postHandler.GetByID)     // GET  /api/posts?id=1
		r.Get("/page", postHandler.List)    // GET  /api/posts/page?offset=0&limit=10
	})
}
