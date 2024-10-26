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
	goalRepo := repositories.NewGoalRepository(db)
	hydrationRepo := repositories.NewHydrationRepository(db)
	mealRepo := repositories.NewMealRepository(db)
	sleepRepo := repositories.NewSleepRepository(db)

	// Initialize services
	exerciseService := services.NewExerciseService(exerciseRepo)
	userService := services.NewUserService(userRepo)
	goalService := services.NewGoalService(goalRepo)
	hydrationService := services.NewHydrationService(hydrationRepo)
	sleepService := services.NewSleepService(sleepRepo)
	mealService := services.NewMealService(mealRepo)

	// Initialize controllers
	exerciseController := controllers.NewExerciseController(exerciseService)
	userController := controllers.NewUserController(userService, jwtService)
	goalController := controllers.NewGoalController(goalService)
	hydrationController := controllers.NewHydrationController(hydrationService)
	sleepController := controllers.NewSleepController(sleepService)
	mealController := controllers.NewMealController(mealService)

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
				exercises.GET("/", exerciseController.GetUserExercises)
				exercises.PUT("/:id", exerciseController.Update)
				exercises.DELETE("/:id", exerciseController.Delete)
			}

			// Goal routes
			goals := protected.Group("/goals")
			{
				goals.POST("", goalController.Create)
				//goals.GET("", controllers.Goal.GetAll)
				//goals.GET("/active", controllers.Goal.GetActiveGoals)
				//goals.GET("/:id", controllers.Goal.GetByID)
				//goals.PUT("/:id", controllers.Goal.Update)
				//goals.DELETE("/:id", controllers.Goal.Delete)
				//goals.GET("/:id/progress", controllers.Goal.GetProgress)
			}

			// Hydration routes
			hydration := protected.Group("/hydration")
			{
				hydration.POST("", hydrationController.Create)
				//hydration.GET("", controllers.Hydration.GetAll)
				//hydration.GET("/daily", controllers.Hydration.GetDailySummary)
				//hydration.PUT("/:id", controllers.Hydration.Update)
				//hydration.DELETE("/:id", controllers.Hydration.Delete)
			}

			// Meal routes
			meals := protected.Group("/meals")
			{
				meals.POST("", mealController.Create)
				//meals.GET("", controllers.Meal.GetAll)
				//meals.GET("/:id", controllers.Meal.GetByID)
				//meals.PUT("/:id", controllers.Meal.Update)
				//meals.DELETE("/:id", controllers.Meal.Delete)
				//meals.GET("/summary", controllers.Meal.GetDailySummary)
			}

			// Sleep routes
			sleep := protected.Group("/sleep")
			{
				sleep.POST("", sleepController.Create)
				//sleep.GET("", controllers.Sleep.GetAll)
				//sleep.GET("/weekly", controllers.Sleep.GetWeeklySummary)
				//sleep.PUT("/:id", controllers.Sleep.Update)
				//sleep.DELETE("/:id", controllers.Sleep.Delete)
			}
		}
	}

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Print(port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
