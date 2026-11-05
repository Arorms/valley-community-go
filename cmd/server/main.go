package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/config"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/handler"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/router"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/view"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := repository.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	oidcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oidcClient, err := auth.NewOIDCClient(oidcCtx, cfg.OIDC, cfg.Server.BaseURL)
	if err != nil {
		log.Fatalf("OIDC client init failed: %v", err)
	}

	secureCookie := strings.HasPrefix(cfg.Server.BaseURL, "https://")
	sessions := auth.NewSessionManager(cfg.Session.Secret, secureCookie)

	renderer, err := view.NewRenderer()
	if err != nil {
		log.Fatalf("Template init failed: %v", err)
	}

	// repositories
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	profileRepo := repository.NewUserProfileRepository(db)

	// services
	postService := service.NewPostService(postRepo, commentRepo, likeRepo, favoriteRepo, profileRepo)
	commentService := service.NewCommentService(db, postRepo, commentRepo, profileRepo, notificationRepo)
	likeService := service.NewLikeService(db, postRepo, commentRepo, likeRepo, notificationRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, postRepo, likeRepo, profileRepo)
	notificationService := service.NewNotificationService(notificationRepo, postRepo, commentRepo, profileRepo)
	profileService := service.NewUserProfileService(profileRepo)

	base := &handler.Base{
		Renderer:      renderer,
		Sessions:      sessions,
		Posts:         postService,
		Comments:      commentService,
		Likes:         likeService,
		Favorites:     favoriteService,
		Notifications: notificationService,
		Profiles:      profileService,
	}

	r := router.New(router.Handlers{
		Web: router.WebHandlers{
			Home:         handler.NewHomeHandler(base),
			Post:         handler.NewPostHandler(base),
			Comment:      handler.NewCommentHandler(base),
			Like:         handler.NewLikeHandler(base),
			Favorite:     handler.NewFavoriteHandler(base),
			Notification: handler.NewNotificationHandler(base),
			Profile:      handler.NewProfileHandler(base),
			Auth:         handler.NewAuthHandler(base, oidcClient),
		},
	}, sessions)

	log.Printf("Server is running on %s", cfg.Server.BaseURL)
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
