package main

import (
	"fmt"
	"log"

	"dict-hub/internal/config"
	"dict-hub/internal/database"
	"dict-hub/internal/router"
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

	// 如果配置了自动加载，则加载字典目录
	if cfg.MDX.AutoLoad && cfg.MDX.DictDir != "" {
		log.Printf("Auto-loading dictionaries from: %s", cfg.MDX.DictDir)
		if err := mdxManager.LoadAll(cfg.MDX.DictDir); err != nil {
			log.Printf("Warning: Failed to auto-load dictionaries: %v", err)
		} else {
			dicts := mdxManager.ListLoaded()
			log.Printf("Loaded %d dictionaries", len(dicts))
		}
	}

	r := router.Setup(cfg, db, mdxManager)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("MDX API: http://localhost%s/api/v1/dicts", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
