package v1handler

import (
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
	return &UsersHandler{
		service: service,
	}
}

func (uh *UsersHandler) GetUserByUuid(ctx *gin.Context) {
	idStr := ctx.Param("id")
	// Kiểm tra định dạng uuid
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err.Error())
		return
	}

	user, err := uh.service.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving user", err.Error())
		return
	}

	// Dùng SuccessResponse cho đồng bộ với các hàm khác
	utils.SuccessResponse(ctx, http.StatusOK, "User retrieved successfully", user)
}

func (uh *UsersHandler) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "CreateUser called",
	})
}

func (uh *UsersHandler) GetProfile(c *gin.Context) {
	// 1. Lấy userID từ Context (Do AuthMiddleware đã ghim vào)
	userIDStr, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized", "Không tìm thấy thông tin người dùng")
		return
	}
	// 2. Parse ID sang UUID
	id, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, 400, "Bad Request", "ID người dùng không hợp lệ")
		return
	}
	// 3. Gọi Repo lấy thông tin user từ DB
	user, err := uh.service.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, 404, "Not Found", "Người dùng không tồn tại")
		return
	}
	// 4. Trả về thông tin (Password sẽ tự ẩn do json:"-" trong model)
	utils.SuccessResponse(c, 200, "Lấy thông tin cá nhân thành công", user)
}

// Hàm GetAll (Lấy danh sách có phân trang)
func (uh *UsersHandler) GetAll(c *gin.Context) {
	//Lấy page và limit từ URL (Query Param). Ví dụ: ?page=1&limit=10.
	//DefaultQuery khi người dùng ko nhap gì cả mạc dịnh sẽ là nó
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	// 2. Convert string -> int
	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)
	// Nếu khách nhập tào lao
	if errPage != nil || errLimit != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid parameters", "Page and limit must be numbers")
		return
	}
	//3 gọi service
	//users ds user lấy được , total tông sl user có trong data
	users, total, err := uh.service.GetAllUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Server Error", err.Error())
		return
	}
	// 4. Trả về format chuẩn có cả thông tin phân trang
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users successfully",
		"data":    users,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Hàm Update (Cập nhật thông tin)
func (uh *UsersHandler) Update(c *gin.Context) {
	// 1. Lấy ID từ URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", err.Error())
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// ShouldBindJSON tự động check các rules validation nếu có
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Input", err.Error())
		return
	}
	// 3. Gọi Service xử lý logic update
	if err := uh.service.UpdateUser(id, req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Update Failed", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

func (uh *UsersHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", err.Error())
		return
	}
	if err := uh.service.DeleteUser(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Delete Failed", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

func (uh *UsersHandler) Search(c *gin.Context) {
	// Lấy keyword từ URL: /users/search?keyword=Hung
	keyword := c.Query("keyword")
	users, err := uh.service.SearchUsers(keyword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Search Failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Search completed", users)
}
