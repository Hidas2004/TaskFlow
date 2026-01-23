package middlewares

import (
	"errors" // Cần thêm cái này để tạo lỗi từ string
	"net/http"
	"strings"

	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// BƯỚC 1: LẤY TOKEN
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {

			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("token missing or invalid format"))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// BƯỚC 2: VALIDATE TOKEN
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {

			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid Token", err)
			return
		}

		// BƯỚC 3: Ép kiểu UserID từ String sang UUID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid User ID in Token", err)
			return
		}

		// Lưu UUID chuẩn vào context
		c.Set("userID", userID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// RoleMiddleware kiểm tra quyền hạn
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy thông tin role đã lưu từ AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("role info not found in context"))
			return
		}

		// 2. Ép kiểu dữ liệu từ interface{} sang string
		roleStr := userRole.(string)

		// 3. Kiểm tra xem role của user có nằm trong danh sách cho phép không
		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next() // Có quyền -> Đi tiếp
				return
			}
		}

		// 4. Nếu không khớp quyền nào -> Chặn
		utils.ErrorResponse(c, http.StatusForbidden, "Access Denied", errors.New("insufficient permissions"))
	}
}
