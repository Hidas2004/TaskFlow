package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.RouterGroup, taskHandler *v1handler.TaskHandler, authMiddleware gin.HandlerFunc) {
	tasks := router.Group("/tasks")
	tasks.Use(authMiddleware) // Bắt buộc đăng nhập
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("", taskHandler.GetTasks)        // Đã gộp Search, Filter, GetMyTasks
		tasks.GET("/:id", taskHandler.GetTaskByID) // Bạn tự code hàm này nhé, tương tự UpdateTask
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}
