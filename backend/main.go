package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"valley-community-go/handlers"
	"valley-community-go/repositories"
	"valley-community-go/routes"
	"valley-community-go/services"
)

func main() {
	log.SetOutput(os.Stdout)

	db, err := repositories.Connect()
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	// DI
	postRepository := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService)

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

	routes.RegisterPostRoutes(r, postHandler)

	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
