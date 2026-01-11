package main

import (
	"fmt"
	"log"

	"dict-hub/internal/config"
	"dict-hub/internal/database"
	"dict-hub/internal/router"
	"dict-hub/internal/service"
	"dict-hub/internal/service/mdx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}
	defer sqlDB.Close()

	// 初始化 MDX 字典管理器
	mdxManager := mdx.NewManager()

	// 初始化服务
	historySvc := service.NewHistoryService(db)
	wordFreqSvc := service.NewWordFreqService(db)
	dictSourceSvc := service.NewDictSourceService(db, mdxManager, cfg.MDX.DictDir)
	downloadSvc := service.NewDownloadService(db, cfg.MDX.DictDir, dictSourceSvc)

	// 启动时同步：从数据库加载已启用的字典
	log.Printf("Syncing dictionaries from database...")
	if err := dictSourceSvc.SyncOnStartup(); err != nil {
		log.Printf("Warning: Failed to sync dictionaries: %v", err)
	} else {
		dicts := mdxManager.ListLoaded()
		log.Printf("Synced %d dictionaries from database", len(dicts))
	}

	// 组装服务
	svcs := &router.Services{
		DictSourceSvc: dictSourceSvc,
		DownloadSvc:   downloadSvc,
		HistorySvc:    historySvc,
		WordFreqSvc:   wordFreqSvc,
	}

	r := router.Setup(cfg, db, mdxManager, svcs)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("MDX API: http://localhost%s/api/v1/dicts", addr)
	log.Printf("Dictionary Management API: http://localhost%s/api/v1/dictionaries", addr)
	log.Printf("History API: http://localhost%s/api/v1/history", addr)
	log.Printf("Word Frequency API: http://localhost%s/api/v1/wordfreq", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
