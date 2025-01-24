package main

import (
	"golang/internal/config"
	"golang/internal/handlers"
	"golang/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	app := gin.Default()

	err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading environment variables")
		return
	}

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to my public API"})
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
