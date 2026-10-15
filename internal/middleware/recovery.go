package middleware

import (
	"log"
	"net/http"

	"booking-studio-foto/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] %v", r)
				response.Error(c, http.StatusInternalServerError, "Internal server error occurred", gin.H{
					"error": "An unexpected error occurred. Please try again later.",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
