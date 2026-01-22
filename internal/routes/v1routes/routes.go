package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRoutes: Nhận tất cả Handler và cấu hình router
func SetupRoutes(router *gin.Engine, cfg *config.Config,
	authHandler *v1handler.AuthHandler,
	usersHandler *v1handler.UsersHandler,
	teamHandler *v1handler.TeamHandler, // <--- Thêm mới
	taskHandler *v1handler.TaskHandler, // <--- Thêm mới
	commentHandler *v1handler.CommentHandler,
	attachmentHandler *v1handler.AttachmentHandler, // <--- Thêm mới
) {

	api := router.Group("/api/v1")

	// 1. Route Public (Không cần Login)
	SetupAuthRoutes(api, authHandler)

	// 2. Route Protected (Cần Login)
	// Middleware được gắn chung ở đây 1 lần duy nhất, các route con không cần truyền nữa
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	{
		SetupUserRoutes(protected, usersHandler)
		SetupCommentRoutes(protected, commentHandler)

		// Gọi hàm đăng ký cho Task (Hàm này chúng ta đã sửa ở bước trước)
		RegisterTaskRoutes(protected, taskHandler)

	}

	RegisterTeamRoutes(protected, cfg, teamHandler)
	RegisterAttachmentRoutes(protected, attachmentHandler) // Bỏ authMiddleware vì đã có ở Group cha
}
