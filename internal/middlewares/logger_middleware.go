package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"healthApi/internal/utils"
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
