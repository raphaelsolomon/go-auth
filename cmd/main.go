package main

import (
	"golang/internal/config"
	"golang/internal/database"
	"golang/internal/handlers"
	"golang/internal/middleware"
	"golang/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	app := gin.Default()

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading environment variables")
		return
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	// Auto-migrate database models
	if err := database.GetDB().AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to my public API"})
	})

	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	//Auth routes
	auth := app.Group("/api/auth")
	auth.POST("/register", handlers.RegisterHanders)
	auth.POST("/login", handlers.LoginHanders)

	protected := app.Group("/api/v1")
	protected.Use(middleware.AuthMiddleWare())

	protected.GET("/user", handlers.GetUserHandlers)
	protected.PUT("/user", handlers.UpdateUserHandlers)

	port := config.GetEnv("PORT")
	log.Println("Starting server on port", port)
	app.Run(":" + port)
}
