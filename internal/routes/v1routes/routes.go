package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRoutes: Hàm Tổng quản lý Route được gọi ở main.go
func SetupRoutes(router *gin.Engine, cfg *config.Config, authHandler *v1handler.AuthHandler, usersHandler *v1handler.UsersHandler) {

	api := router.Group("/api/v1")

	SetupAuthRoutes(api, authHandler)

	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	SetupUserRoutes(protected, usersHandler)

}
