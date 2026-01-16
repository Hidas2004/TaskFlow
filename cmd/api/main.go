package main

import (
	"log"

	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/Hidas2004/TaskFlow/internal/routes/v1routes"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load cấu hình
	cfg := config.LoadConfig()

	// 2. Kết nối DB
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối: %v", err)
	}

	// 3. Migration (Tạo bảng tự động)
	// Phase 5, 6, 7 sẽ thêm các model khác vào đây (Task, Team...)
	if err := db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
		&models.Task{},
		&models.Comment{},
	); err != nil {
		log.Fatalf("❌ Lỗi migration: %v", err)
	}
	log.Println("✅ Migration thành công!")

	// 4. Khởi tạo Gin Router
	router := gin.Default()

	// Tầng 1: Repository (Giao tiếp DB)
	userRepo := repositories.NewUserRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	// Tầng 2: Service (Xử lý logic, cần Repo và Config)
	userService := v1services.NewUserService(userRepo)
	authService := v1services.NewAuthService(userRepo, cfg)
	teamService := v1services.NewTeamService(teamRepo, userRepo)
	taskService := v1services.NewTaskService(taskRepo, teamRepo)

	// Tầng 3: Handler (Xử lý HTTP, cần Service)
	authHandler := v1handler.NewAuthHandler(authService)
	teamHandler := v1handler.NewTeamHandler(teamService)
	usersHandler := v1handler.NewUsersHandler(userService)
	taskHandler := v1handler.NewTaskHandler(taskService)

	//Middleware
	authMiddleware := middlewares.AuthMiddleware(cfg.JWTSecret)

	// ==========================================
	// 6. SETUP ROUTES (Cấu hình đường dẫn)
	// ==========================================

	apiV1 := router.Group("/api/v1")

	// Gọi hàm "Tổng quản" SetupRoutes từ package v1routes
	// Hàm này sẽ tự chia route Public và Protected (có Middleware)
	v1routes.SetupRoutes(router, cfg, authHandler, usersHandler)
	v1routes.RegisterTeamRoutes(apiV1, cfg, teamHandler)
	v1routes.RegisterTaskRoutes(apiV1, taskHandler, authMiddleware)
	// 7. Chạy server
	log.Printf("🚀 Server đang chạy tại cổng: %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Không thể khởi động server: %v", err)
	}
}
