package main

import (
	"log"

	"github.com/Hidas2004/TaskFlow/internal/app"
	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load cấu hình
	cfg := config.LoadConfig()

	// 2. Kết nối db
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối: %v", err)
	}

	// 3. Khởi tạo gin router
	router := gin.Default()

	// 4. Khởi tạo Context và Module
	moduleCtx := app.NewModuleContext(db, cfg)
	authModule := app.NewAuthModule(moduleCtx)

	// 5. Đăng ký routes
	apiGroup := router.Group("/api/v1")
	authModule.GetRoutes().Register(apiGroup)

	// 6. Chạy server
	log.Printf("Server running on port %s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
