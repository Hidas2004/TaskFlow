package v1handler

import (
	"errors" // <--- BẮT BUỘC PHẢI CÓ
	"net/http"
	"strconv"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersHandler struct {
	service v1services.UserService
}

func NewUsersHandler(service v1services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

// Helper: Lấy UserID an toàn (Tránh panic)
func (uh *UsersHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	idVal, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	// Case 1: Middleware lưu UUID object
	if id, ok := idVal.(uuid.UUID); ok {
		return id, true
	}
	// Case 2: Middleware lưu string
	if idStr, ok := idVal.(string); ok {
		if id, err := uuid.Parse(idStr); err == nil {
			return id, true
		}
	}
	return uuid.Nil, false
}

func (uh *UsersHandler) GetProfile(c *gin.Context) {
	// 1. Lấy ID an toàn
	userID, ok := uh.getUserIDFromContext(c)
	if !ok {
		// [ĐÃ SỬA]: Dùng errors.New(...)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("user ID invalid or missing"))
		return
	}

	// 2. Gọi Service
	user, err := uh.service.GetUserByID(userID)
	if err != nil {
		utils.HandleServiceError(c, err) // Tự trả về 404 nếu không thấy
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Get profile success", user)
}

func (uh *UsersHandler) GetUserByUuid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	user, err := uh.service.GetUserByID(id)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Get user details success", user)
}

func (uh *UsersHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)

	if errPage != nil || errLimit != nil {
		// [ĐÃ SỬA]: Dùng errors.New(...)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid parameters", errors.New("page and limit must be numbers"))
		return
	}

	users, total, err := uh.service.GetAllUsers(page, limit)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Dùng Helper Pagination để chuẩn hóa output giống TaskHandler
	response := utils.NewPaginationResponse(users, page, limit, total)
	utils.SuccessResponse(c, http.StatusOK, "Get user list success", response)
}

func (uh *UsersHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", err)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Input", err)
		return
	}

	if err := uh.service.UpdateUser(id, req); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

func (uh *UsersHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", err)
		return
	}

	if err := uh.service.DeleteUser(id); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

func (uh *UsersHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	users, err := uh.service.SearchUsers(keyword)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Search users success", users)
}

func (uh *UsersHandler) CreateUser(c *gin.Context) {
	// Placeholder function
	utils.SuccessResponse(c, http.StatusOK, "CreateUser called", nil)
}
