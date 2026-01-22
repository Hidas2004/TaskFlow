package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

func SetupCommentRoutes(router *gin.RouterGroup, handler *v1handler.CommentHandler) {
	comments := router.Group("/comments")
	{
		comments.PUT("/:id", handler.UpdateComment)
		comments.DELETE("/:id", handler.DeleteComment)
	}

	tasks := router.Group("/tasks")
	{
		// SỬA Ở ĐÂY: Đổi :taskId -> :id
		tasks.POST("/:id/comments", handler.CreateComment)
		tasks.GET("/:id/comments", handler.GetComments)
	}
}
