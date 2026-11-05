package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/handler"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/middleware"
)

type Handlers struct {
	Post *handler.PostHandler
	Auth *handler.AuthHandler
}

func New(handlers Handlers, tokenValidator *middleware.TokenValidator) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)

	r.Get("/api/post", handlers.Post.GetByID)
	r.Get("/api/post/page", handlers.Post.List)

	r.Group(func(r chi.Router) {
		r.Use(tokenValidator.Authenticate)
		r.Post("/api/post", handlers.Post.CreatePost)
		r.Get("/api/auth/me", handlers.Auth.Me)
	})

	return r
}
