package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CorsMiddleware() gin.HandlerFunc {
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true, // Enable Debugging for testing, consider disabling in production
	})

	return func(c *gin.Context) {
		corsConfig.HandlerFunc(c.Writer, c.Request)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
