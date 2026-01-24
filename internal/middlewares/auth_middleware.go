package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware: Kiểm tra Token và xác thực user
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- DEBUG LOG START (Dùng để bắt lỗi 401) ---
		authHeader := c.GetHeader("Authorization")
		fmt.Println("-------------------------------------------")
		fmt.Println("[DEBUG] 1. Auth Header nhận được:", authHeader)
		// --- DEBUG LOG END ---

		// BƯỚC 1: Kiểm tra format Header
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("[DEBUG] ❌ Lỗi: Header thiếu hoặc sai format")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("token missing or invalid format"))
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// BƯỚC 2: Validate Token
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			fmt.Println("[DEBUG] ❌ Lỗi Validate Token:", err.Error())
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid Token", err)
			c.Abort()
			return
		}

		// BƯỚC 3: Parse UserID
		fmt.Println("[DEBUG] 3. Claims UserID:", claims.UserID)
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			fmt.Println("[DEBUG] ❌ Lỗi Parse UUID:", err.Error())
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid User ID in Token", err)
			c.Abort()
			return
		}

		// BƯỚC 4: Lưu thông tin vào Context để dùng ở các handler sau
		c.Set("userID", userID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		fmt.Println("[DEBUG] ✅ Auth Success! UserID:", userID)
		c.Next()
	}
}

// RoleMiddleware: Kiểm tra quyền hạn (Admin, Team Leader, Member...)
// Hàm này phải nằm CÙNG FILE với AuthMiddleware nhưng ở bên ngoài
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy thông tin role đã lưu từ AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("role info not found in context"))
			c.Abort()
			return
		}

		// 2. Ép kiểu dữ liệu từ interface{} sang string
		roleStr, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", errors.New("invalid role format"))
			c.Abort()
			return
		}

		// 3. Kiểm tra xem role của user có nằm trong danh sách cho phép không
		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next() // Có quyền -> Đi tiếp
				return
			}
		}

		// 4. Nếu không khớp quyền nào -> Chặn
		utils.ErrorResponse(c, http.StatusForbidden, "Access Denied", errors.New("insufficient permissions"))
		c.Abort()
	}
}
