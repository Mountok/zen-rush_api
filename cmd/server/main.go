package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zenrush/backend/internal/db"
	"github.com/zenrush/backend/internal/handlers"
	"github.com/zenrush/backend/internal/middleware"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	r := gin.Default()

	// Настройка CORS для фронтенда
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5500", "http://localhost:5173", "http://localhost:3000", "http://localhost:4173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)

		activities := api.Group("/activities")
		activities.Use(middleware.JWTAuth())
		activities.GET("", handlers.ListActivities)
		activities.POST("", handlers.CreateActivity)
		activities.GET(":id", handlers.GetActivity)
		activities.PUT(":id", handlers.UpdateActivity)
		activities.DELETE(":id", handlers.DeleteActivity)

		favorites := api.Group("/favorites")
		favorites.Use(middleware.JWTAuth())
		favorites.GET("", handlers.ListFavorites)
		favorites.POST(":activity_id", handlers.AddFavorite)
		favorites.DELETE(":activity_id", handlers.RemoveFavorite)

		history := api.Group("/history")
		history.Use(middleware.JWTAuth())
		history.GET("", handlers.ListHistory)
		history.POST(":activity_id", handlers.AddHistory)

		// --- Mood stats ---
		api.POST("/mood-stats", middleware.JWTAuth(), handlers.SaveOrUpdateMoodStat)
		api.GET("/users/me/mood-stats", middleware.JWTAuth(), handlers.GetMoodStats)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
