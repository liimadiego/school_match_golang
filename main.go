package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/liimadiego/schoolmatch/internal/database"
	"github.com/liimadiego/schoolmatch/internal/handlers"
	"github.com/liimadiego/schoolmatch/internal/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main() {
	database.Connect()
	database.Migrate()

	r := gin.Default()

	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)
	r.GET("/api/schools", handlers.GetSchools)
	r.GET("/api/schools/:id", handlers.GetSchool)
	r.GET("/api/school-reviews/:school_id", handlers.GetReviews)
	r.GET("/api/reviews/:id", handlers.GetReview)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/schools", handlers.CreateSchool)
		auth.PUT("/schools/:id", handlers.UpdateSchool)
		auth.DELETE("/schools/:id", handlers.DeleteSchool)

		auth.POST("/reviews", handlers.CreateReview)
		auth.PUT("/reviews/:id", handlers.UpdateReview)
		auth.DELETE("/reviews/:id", handlers.DeleteReview)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
