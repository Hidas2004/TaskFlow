package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

func RegisterAttachmentRoutes(router *gin.RouterGroup, handler *v1handler.AttachmentHandler) {

	tasks := router.Group("/tasks")
	{
		tasks.POST("/:id/attachments", handler.Upload)
		tasks.GET("/:id/attachments", handler.GetByTask)
	}

	attachments := router.Group("/attachments")
	{
		attachments.DELETE("/:id", handler.Delete)
	}
}
