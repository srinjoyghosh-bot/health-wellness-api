package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"healthApi/internal/config"
	"healthApi/internal/controllers"
	"healthApi/internal/middlewares"
	"healthApi/internal/repositories"
	"healthApi/internal/services"
	"healthApi/internal/utils"
	"healthApi/pkg/database"
	"log"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.InitDB()

	// Initialize JWT service
	jwtService := utils.NewJWTService(
		viper.GetString("jwt.secret"),
		viper.GetDuration("jwt.expiration")*time.Hour,
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	exerciseRepo := repositories.NewExerciseRepository(db)

	// Initialize services
	exerciseService := services.NewExerciseService(exerciseRepo)
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	exerciseController := controllers.NewExerciseController(exerciseService)
	userController := controllers.NewUserController(userService, jwtService)

	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middlewares.Logger())

	// Routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware(jwtService))
		{
			// Exercise routes
			exercises := protected.Group("/exercises")
			{
				exercises.POST("/", exerciseController.CreateExercise)
				exercises.GET("/:id", exerciseController.GetUserExercises)
				exercises.PUT("/:id", exerciseController.Update)
				exercises.DELETE("/:id", exerciseController.Delete)
			}

			// Similar route groups for meals, sleep, hydration, and goals...
		}
	}

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Print(port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
