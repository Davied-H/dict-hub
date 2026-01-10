package router

import (
	"dict-hub/internal/config"
	"dict-hub/internal/handler"
	"dict-hub/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		_ = api
	}

	return r
}
