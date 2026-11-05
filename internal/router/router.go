package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/handler"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/view"
)

type WebHandlers struct {
	Home         *handler.HomeHandler
	Post         *handler.PostHandler
	Comment      *handler.CommentHandler
	Like         *handler.LikeHandler
	Favorite     *handler.FavoriteHandler
	Notification *handler.NotificationHandler
	Profile      *handler.ProfileHandler
	Auth         *handler.AuthHandler
}

type Handlers struct {
	Web WebHandlers
}

func New(handlers Handlers, sessions *auth.SessionManager) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(sessions.LoadUser)

	staticFiles := http.FileServer(http.FS(view.StaticFS()))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFiles))

	// 公开页面
	r.Get("/", handlers.Web.Home.Index)
	r.Get("/posts/{id}", handlers.Web.Post.Show)
	r.Get("/users/{username}", handlers.Web.Profile.Show)
	r.Get("/auth/login", handlers.Web.Auth.Login)
	r.Get("/auth/callback", handlers.Web.Auth.Callback)

	// 需要登录的页面与表单
	r.Group(func(pr chi.Router) {
		pr.Use(auth.RequireAuth)
		pr.Use(sessions.CSRF)

		pr.Get("/posts/new", handlers.Web.Post.New)
		pr.Post("/posts", handlers.Web.Post.Create)
		pr.Post("/posts/{id}/comments", handlers.Web.Comment.Create)
		pr.Post("/posts/{id}/like", handlers.Web.Like.TogglePost)
		pr.Post("/comments/{id}/like", handlers.Web.Like.ToggleComment)
		pr.Post("/posts/{id}/favorite", handlers.Web.Favorite.Toggle)
		pr.Get("/favorites", handlers.Web.Favorite.List)
		pr.Get("/notifications", handlers.Web.Notification.Index)
		pr.Post("/notifications/read", handlers.Web.Notification.MarkAllRead)
		pr.Post("/auth/logout", handlers.Web.Auth.Logout)
	})

	return r
}
