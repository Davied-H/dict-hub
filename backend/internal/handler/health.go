package handler

import (
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "healthy",
		"service": "dict-hub",
		"version": "1.0.0",
	})
}
