package main

import (
	"log"
	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// loading env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// loading config
	cfg := config.Load()

	// connect to database
	if err := repository.ConnectDB(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// create router
	router := gin.Default()

	// create handlers
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	todoHandler := handlers.NewTodoHandler()

	// public routes
	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	// authorized routes group
	authGroup := router.Group("/api")
	authGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authGroup.GET("/todos", todoHandler.GetTodos)
		authGroup.POST("/todos", todoHandler.CreateTodo)
		authGroup.PUT("/todos/:id", todoHandler.UpdateTodo)
		authGroup.DELETE("/todos/:id", todoHandler.DeleteTodo)
	}

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "todo-api",
			"version": "1.0.0",
		})
	})

	// run server
	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
