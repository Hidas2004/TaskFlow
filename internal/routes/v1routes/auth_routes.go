// Định tuyến (Routes)
// AuthRoutes (File 2): Là lễ tân. Họ cầm bản đồ bàn,
//
//	hướng dẫn khách đi vào đúng chỗ. Khi khách muốn
//	"Login", lễ tân sẽ dẫn khách đến gặp đầu bếp AuthHandler
package v1routes

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/gin-gonic/gin"
)

// : Struct AuthRoutes không tự xử lý logic,
// nó chứa (giữ tham chiếu đến) AuthHandler
type AuthRoutes struct {
	handler *v1handler.AuthHandler
}

// (handler *v1handler.AuthHandler) Đây là tham số đầu vào
//
// Hàm này yêu cầu một con trỏ (*) trỏ tới struct
//
// AuthHandler (nằm trong package v1handler).
func NewAuthRoutes(handler *v1handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{handler: handler}
}

func (ar *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", ar.handler.Register)
		auth.POST("/login", ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
	}

}
