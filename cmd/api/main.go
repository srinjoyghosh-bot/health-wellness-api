package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"healthApi/internal/config"
	"healthApi/internal/controllers"
	"healthApi/internal/repositories"
	"healthApi/internal/services"
	"healthApi/pkg/database"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.InitDB()

	// Initialize repositories
	exerciseRepo := repositories.NewExerciseRepository(db)

	// Initialize services
	exerciseService := services.NewExerciseService(exerciseRepo)

	// Initialize controllers
	exerciseController := controllers.NewExerciseController(exerciseService)

	router := gin.Default()

	// Routes
	api := router.Group("/api")
	{
		// Exercise routes
		exercises := api.Group("/exercises")
		{
			exercises.POST("/", exerciseController.CreateExercise)
			exercises.GET("/", exerciseController.GetUserExercises)
		}

		// Add other route groups...
	}

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Print(port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
