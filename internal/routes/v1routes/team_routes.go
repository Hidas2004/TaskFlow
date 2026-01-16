package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterTeamRoutes(router *gin.RouterGroup, cfg *config.Config, teamHandler *v1handler.TeamHandler) {

	teams := router.Group("/teams")

	teams.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	{
		teams.POST("", teamHandler.Create)
		teams.GET("", teamHandler.GetAll)
		teams.GET("/my", teamHandler.GetMyTeams)
		teams.GET("/:id", teamHandler.GetByID)
		teams.PUT("/:id", teamHandler.Update)
		teams.DELETE("/:id", teamHandler.Delete)
		teams.POST("/:id/members", teamHandler.AddMember)
		teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
		teams.GET("/:id/members", teamHandler.GetMembers)
	}
}
