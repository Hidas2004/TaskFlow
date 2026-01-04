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
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Dữ liệu không hợp lệ", err.Error())
		return
	}

	// Gọi service với input là struct DTO (QUAN TRỌNG: Đã sửa chỗ này)
	resp, err := ah.service.Login(&input)
	if err != nil {
		utils.ErrorResponse(ctx, 401, "Đăng nhập thất bại", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "Đăng nhập thành công", resp)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {

}

// 1. Xử lý Register
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var input dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Dữ liệu không hợp lệ", err.Error())
		return
	}

	resp, err := ah.service.Register(&input) // Truyền pointer struct
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Đăng ký thất bại", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "Đăng ký thành công", resp)
}
