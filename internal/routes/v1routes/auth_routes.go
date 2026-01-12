package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes: Chỉ đăng ký các đường dẫn liên quan đến Auth
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *v1handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
