package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/config"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/handler"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/middleware"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/router"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
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

	// NewTokenValidator 会拉取 OIDC discovery 文档，给个超时避免启动卡死
	oidcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokenValidator, err := middleware.NewTokenValidator(oidcCtx, middleware.OIDCOptions{
		IssuerURL: cfg.OIDC.IssuerURL,
		Audience:  cfg.OIDC.Audience,
		ClientID:  cfg.OIDC.ClientID,
	})
	if err != nil {
		log.Fatalf("OIDC init failed: %v", err)
	}

	// DI
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler()

	// Register routes
	r := router.New(router.Handlers{
		Post: postHandler,
		Auth: authHandler,
	}, tokenValidator)

	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
