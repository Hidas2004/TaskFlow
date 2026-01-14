package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterTeamRoutes(router *gin.Engine, cfg *config.Config, teamHandler *v1handler.TeamHandler) {
	//1 tạo đương dẫn
	teams := router.Group("/api/teams")
	//2 cài bảo mất middelware
	// bát kỳ ái muốn vào thì phải có token phu hợp
	teams.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
	// 3. Định nghĩa các endpoints cụ thể (Đúng như đoạn code bạn gửi)
	{
		teams.POST("", teamHandler.Create)                             // Tạo team
		teams.GET("", teamHandler.GetAll)                              // Lấy hết team
		teams.GET("/my", teamHandler.GetMyTeams)                       // Lấy team của tôi
		teams.GET("/:id", teamHandler.GetByID)                         // Lấy chi tiết team
		teams.PUT("/:id", teamHandler.Update)                          // Sửa team
		teams.DELETE("/:id", teamHandler.Delete)                       // Xóa team
		teams.POST("/:id/members", teamHandler.AddMember)              // Thêm thành viên
		teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember) // Xóa thành viên
		teams.GET("/:id/members", teamHandler.GetMembers)              // Xem thành viên
	}
}
