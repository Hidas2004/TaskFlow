package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// Đã thêm tham số 'secretKey string' vào hàm này
func RegisterCommentRoutes(router *gin.RouterGroup, handler *v1handler.CommentHandler, secretKey string) {

	// Group các API cần đăng nhập mới dùng được
	commentGroup := router.Group("/comments")
	// TRUYỀN secretKey VÀO ĐÂY
	commentGroup.Use(middlewares.AuthMiddleware(secretKey))
	{
		// PUT /api/v1/comments/:id -> Sửa
		commentGroup.PUT("/:id", handler.UpdateComment)
		// DELETE /api/v1/comments/:id -> Xóa
		commentGroup.DELETE("/:id", handler.DeleteComment)
	}

	// Riêng phần Task + Comment thường đi chung URL
	taskGroup := router.Group("/tasks")
	// TRUYỀN secretKey VÀO ĐÂY
	taskGroup.Use(middlewares.AuthMiddleware(secretKey))
	{
		// Đã sửa tên hàm GetComments chuẩn chính tả
		taskGroup.POST("/:taskId/comments", handler.CreateComment)
		taskGroup.GET("/:taskId/comments", handler.GetComments)
	}
}
