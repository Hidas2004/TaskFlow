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
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load cấu hình & DB
	cfg := config.LoadConfig()
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối: %v", err)
	}

	utils.RegisterCustomValidators()

	// 2. Migration
	db.AutoMigrate(
		&models.User{}, &models.Team{}, &models.TeamMember{},
		&models.Task{}, &models.Comment{}, &models.Attachment{},
	)

	// 3. Khởi tạo Router
	router := gin.Default()

	router.Use(middlewares.CorsMiddleware(cfg))
	router.Use(middlewares.RateLimitMiddleware(cfg))
	router.Static("/uploads", "./uploads")

	// 4. Khởi tạo Layers (Repo -> Service -> Handler)
	// --- Repo ---
	userRepo := repositories.NewUserRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	attachmentRepo := repositories.NewAttachmentRepository(db)

	// --- Service ---
	userService := v1services.NewUserService(userRepo)
	authService := v1services.NewAuthService(userRepo, cfg)
	teamService := v1services.NewTeamService(teamRepo, userRepo)
	taskService := v1services.NewTaskService(taskRepo, teamRepo)
	commentService := v1services.NewCommentService(commentRepo, userRepo, taskRepo, teamRepo)
	attachmentService := v1services.NewAttachmentService(attachmentRepo, taskRepo, cfg)

	// --- Handler ---
	authHandler := v1handler.NewAuthHandler(authService)
	teamHandler := v1handler.NewTeamHandler(teamService)
	usersHandler := v1handler.NewUsersHandler(userService)
	taskHandler := v1handler.NewTaskHandler(taskService)
	commentHandler := v1handler.NewCommentHandler(commentService)
	attachmentHandler := v1handler.NewAttachmentHandler(attachmentService)

	// 5. SETUP ROUTES (Gọn gàng)
	// Truyền tất cả Handler vào đây, để routes.go tự lo liệu
	v1routes.SetupRoutes(router, cfg,
		authHandler,
		usersHandler,
		teamHandler, // Mới thêm
		taskHandler, // Mới thêm
		commentHandler,
		attachmentHandler, // Mới thêm
	)

	// 6. Chạy Server
	log.Printf("🚀 Server đang chạy tại cổng: %s", cfg.ServerPort)
	log.Printf("🌍 Allowed Origins: %s", cfg.ClientOrigin)
	router.Run(":" + cfg.ServerPort)
}
