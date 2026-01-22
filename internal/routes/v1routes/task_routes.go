package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.RouterGroup, taskHandler *v1handler.TaskHandler) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("/dashboard/stats", taskHandler.GetDashboardStats)

		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("", taskHandler.GetTasks) // Đã gộp Search, Filter, Pagination

		tasks.GET("/:id", taskHandler.GetTaskByID)

		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}
