package middlewares

import (
	"strings"

	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// BƯỚC 1: LẤY TOKEN (Giữ nguyên)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorResponse(c, 401, "Unauthorized", "Token không tồn tại hoặc sai định dạng")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// BƯỚC 2: VALIDATE TOKEN (Giữ nguyên)
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			utils.ErrorResponse(c, 401, "Unauthorized", err.Error())
			c.Abort()
			return
		}

		// --- [SỬA ĐOẠN NÀY] ---
		// BƯỚC 3: Ép kiểu UserID từ String sang UUID ngay tại đây
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			utils.ErrorResponse(c, 401, "Unauthorized", "UserID trong token không hợp lệ")
			c.Abort()
			return
		}

		// Lưu UUID chuẩn vào context (Thay vì lưu string như trước)
		c.Set("userID", userID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// allowedRoles ...string dấu ... nghĩa là bạn có thể truyền 1 hoặc nhiều role tùy
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy thông tin role đã lưu từ AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, 403, "Forbidden", "Không xác định được quyền hạn")
			c.Abort()
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
		utils.ErrorResponse(c, 403, "Forbidden", "Bạn không đủ quyền hạn")
		c.Abort()
	}
}
