package middlewares

import (
	"strings"

	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//BƯỚC 1: LẤY TOKEN
		//token thường được gửi trong header tên là Authorization
		//c.getheader là Lấy giá trị của một HTTP header theo tên
		authHeader := c.GetHeader("Authorization")
		//kiểm tra 1 header có tông tại ko
		//kiểm tra 2 có bắt đầu bằng chữ bearer ko
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorResponse(c, 401, "Unauthorized", "Token không tồn tại hoặc sai định dạng")
			c.Abort()
			return
		}
		// Cắt bỏ chữ "Bearer " (7 ký tự đầu) để lấy chuỗi token sạch
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// --- BƯỚC 2: VALIDATE TOKEN ---
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			utils.ErrorResponse(c, 401, "Unauthorized", err.Error())
			c.Abort()
			return
		}
		// --- BƯỚC 3: LƯU THÔNG TIN (CONTEXT) ---
		// Middleware sau khi kiểm tra vé xong, phải "ghim" thông tin người dùng vào Context.
		// Để lát nữa, các hàm xử lý (Handler) biết được "Ai là người đang gửi request này?".\
		c.Set("userID", claims.UserID)
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
