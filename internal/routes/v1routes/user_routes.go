package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes: Nhận vào router (đã được bọc Middleware ở bên ngoài)
func SetupUserRoutes(router *gin.RouterGroup, userHandler *v1handler.UsersHandler) {

	users := router.Group("/users")
	{

		users.GET("", middlewares.RoleMiddleware("admin"), userHandler.GetAll)

		users.GET("/search", userHandler.Search)
		users.GET("/profile", userHandler.GetProfile)

		users.GET("/:id", userHandler.GetUserByUuid)

		users.PUT("/:id", userHandler.Update)

		// 🛑 Xóa cũng cần khóa, chỉ Admin mới được xóa người khác
		users.DELETE("/:id", middlewares.RoleMiddleware("admin"), userHandler.Delete)
	}
}
