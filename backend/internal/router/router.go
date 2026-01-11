package router

import (
	"dict-hub/internal/config"
	"dict-hub/internal/handler"
	"dict-hub/internal/middleware"
	"dict-hub/internal/service/mdx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, mdxManager mdx.DictManager) *gin.Engine {
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
		// MDX 字典路由
		mdxHandler := handler.NewMdxHandler(mdxManager)
		dicts := api.Group("/dicts")
		{
			dicts.GET("", mdxHandler.List)
			dicts.POST("/load", mdxHandler.Load)
			dicts.POST("/load-all", mdxHandler.LoadAll)
			dicts.GET("/search", mdxHandler.Search)
			dicts.GET("/:id/lookup", mdxHandler.Lookup)
			dicts.GET("/:id/resource/*path", mdxHandler.GetResource)
			dicts.DELETE("/:id", mdxHandler.Unload)
		}
	}

	return r
}
