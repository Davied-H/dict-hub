package main

import (
	"fmt"
	"log"

	"dict-hub/internal/config"
	"dict-hub/internal/database"
	"dict-hub/internal/router"
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

	r := router.Setup(cfg, db)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
