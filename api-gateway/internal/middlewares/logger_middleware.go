package middlewares

import (
	"api-gateway/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var appErr *utils.AppError
			if errors.As(c.Errors[0].Err, &appErr) {
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}
	}
}
