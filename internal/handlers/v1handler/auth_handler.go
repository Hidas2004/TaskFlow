// Logic xử lý (Handler)
// AuthHandler (File 1): Là đầu bếp. Họ biết cách nấu món ăn
// (xử lý Login, Logout) nhưng không tiếp xúc trực tiếp
// với khách ở cửa
package v1handler

import (
	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1services.AuthService
}

// Constructor nhận vào Interface
func NewAuthHandler(service v1services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input dto.LoginRequest
	//kiem tra dau vào
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Dữ liệu không hợp lệ", err.Error())
		return
	}
	//2 goi service de xu ly
	token, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ErrorResponse(ctx, 401, "Đăng nhập thất bại", err.Error())
		return
	}
	//3 trả về kết quả thành công
	utils.SuccessResponse(ctx, 200, "Đăng nhập thành công", gin.H{"token": token})
	
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {

}
