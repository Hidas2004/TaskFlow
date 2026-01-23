package v1handler

import (
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1services.AuthService
}

func NewAuthHandler(service v1services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	resp, err := ah.service.Login(&input)
	if err != nil {

		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", resp)
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	resp, err := ah.service.Register(&input)
	if err != nil {
		utils.HandleServiceError(c, err) // Tự trả về 409 nếu trùng email (nếu service báo lỗi duplicate)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Register successful", resp)
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	// Logout thường xử lý ở client (xóa token), server chỉ cần trả về OK
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
