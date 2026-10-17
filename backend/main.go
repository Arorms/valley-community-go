package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"valley-community-go/config"
	"valley-community-go/handlers"
	"valley-community-go/middlewares"
	"valley-community-go/repositories"
	"valley-community-go/routes"
	"valley-community-go/services"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := repositories.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	// NewTokenValidator 会拉取 OIDC discovery 文档，给个超时避免启动卡死
	oidcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokenValidator, err := middlewares.NewTokenValidator(oidcCtx, middlewares.OIDCOptions{
		IssuerURL: cfg.OIDC.IssuerURL,
		Audience:  cfg.OIDC.Audience,
		ClientID:  cfg.OIDC.ClientID,
	})
	if err != nil {
		log.Fatalf("OIDC init failed: %v", err)
	}

	// DI
	postRepository := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService)
	authHandler := handlers.NewAuthHandler()

	// Register routes
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.RegisterPostRoutes(r, postHandler, tokenValidator)
	routes.RegisterAuthRoutes(r, authHandler, tokenValidator)

	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}