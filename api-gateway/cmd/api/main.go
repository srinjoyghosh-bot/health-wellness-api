package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"healthApi/api-gateway/internal/clients"
	"healthApi/api-gateway/internal/config"
	controllers2 "healthApi/api-gateway/internal/controllers"
	middlewares2 "healthApi/api-gateway/internal/middlewares"
	repositories2 "healthApi/api-gateway/internal/repositories"
	services2 "healthApi/api-gateway/internal/services"
	"healthApi/api-gateway/internal/utils"
	"healthApi/api-gateway/pkg/database"
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
	userRepo := repositories2.NewUserRepository(db)
	exerciseRepo := repositories2.NewExerciseRepository(db)
	goalRepo := repositories2.NewGoalRepository(db)
	hydrationRepo := repositories2.NewHydrationRepository(db)
	mealRepo := repositories2.NewMealRepository(db)
	sleepRepo := repositories2.NewSleepRepository(db)

	// Initialize services
	exerciseService := services2.NewExerciseService(exerciseRepo)
	userService := services2.NewUserService(userRepo)
	goalService := services2.NewGoalService(goalRepo)
	hydrationService := services2.NewHydrationService(hydrationRepo)
	sleepService := services2.NewSleepService(sleepRepo)
	mealService := services2.NewMealService(mealRepo)

	//Initialize gRPC clients
	log.Println("Exercise service address", cfg.Service.ExerciseAddr)
	exerciseClient, err := clients.NewExerciseClient(cfg.Service.ExerciseAddr)
	if err != nil {
		log.Fatalf("Failed to create exercise client: %v", err)
	}
	defer exerciseClient.Close()

	log.Println("User service address", cfg.Service.UserAddr)
	userClient, err := clients.NewUserClient(cfg.Service.UserAddr)
	if err != nil {
		log.Fatalf("Failed to create exercise client: %v", err)
	}
	defer userClient.Close()

	// Initialize controllers
	exerciseController := controllers2.NewExerciseController(exerciseService, exerciseClient)
	userController := controllers2.NewUserController(userService, jwtService, userClient)
	goalController := controllers2.NewGoalController(goalService)
	hydrationController := controllers2.NewHydrationController(hydrationService)
	sleepController := controllers2.NewSleepController(sleepService)
	mealController := controllers2.NewMealController(mealService)

	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middlewares2.Logger())
	router.Use(middlewares2.CorsMiddleware())

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
		protected.Use(middlewares2.AuthMiddleware(jwtService))
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
