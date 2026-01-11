package router

import (
	"time"

	"dict-hub/internal/cache"
	"dict-hub/internal/config"
	"dict-hub/internal/handler"
	"dict-hub/internal/middleware"
	"dict-hub/internal/service"
	"dict-hub/internal/service/mdx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Services 包含所有服务实例
type Services struct {
	DictSourceSvc *service.DictSourceService
	DownloadSvc   *service.DownloadService
	HistorySvc    *service.HistoryService
	WordFreqSvc   *service.WordFreqService
}

func Setup(cfg *config.Config, db *gorm.DB, mdxManager mdx.DictManager, svcs *Services) *gin.Engine {
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	// 创建缓存实例（每分钟清理过期项）
	cacheInstance := cache.New(time.Minute)

	api := r.Group("/api/v1")
	{
		// 搜索路由（新增）
		searchHandler := handler.NewSearchHandler(mdxManager, cacheInstance, db)
		api.GET("/search", searchHandler.Search)
		api.GET("/search/suggest", searchHandler.Suggest)

		// MDX 字典路由（现有，保持兼容）
		mdxHandler := handler.NewMdxHandlerWithHistory(mdxManager, svcs.HistorySvc)

		// 字典查询路由（新增别名）
		api.GET("/dictionaries/:id/lookup", mdxHandler.Lookup)

		// 资源路由（新增别名）
		api.GET("/resources/:dictId/*path", mdxHandler.GetResource)

		// 保留原有 MDX 字典路由（向后兼容）
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

		// 字典管理路由（新增）
		dictHandler := handler.NewDictionaryHandler(svcs.DictSourceSvc, svcs.DownloadSvc)
		dictionaries := api.Group("/dictionaries")
		{
			dictionaries.GET("", dictHandler.List)
			dictionaries.POST("", dictHandler.Add)
			dictionaries.PUT("/:id/toggle", dictHandler.Toggle)
			dictionaries.PUT("/reorder", dictHandler.Reorder)
			dictionaries.POST("/download", dictHandler.Download)
			dictionaries.GET("/download/:taskId", dictHandler.GetDownloadStatus)
			dictionaries.DELETE("/:id", dictHandler.Delete)
		}

		// 历史记录路由（新增）
		historyHandler := handler.NewHistoryHandler(svcs.HistorySvc)
		history := api.Group("/history")
		{
			history.GET("", historyHandler.List)
			history.DELETE("", historyHandler.Clear)
			history.GET("/export", historyHandler.Export)
		}

		// 词频路由（新增）
		wordFreqHandler := handler.NewWordFreqHandler(svcs.WordFreqSvc)
		api.POST("/wordfreq/import", wordFreqHandler.Import)
	}

	return r
}
