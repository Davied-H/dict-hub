package middleware

import (
	"log"
	"net/http"

	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				response.Error(c, http.StatusInternalServerError, 500, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
