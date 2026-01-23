package middlewares

import (
	"strings"
	"time"

	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	if cfg.ClientOrigin == "*" { //Kiểm tra xem trong file .env, bạn có để CLIENT_ORIGIN="*"
		corsConfig.AllowAllOrigins = true
	} else {
		// TRƯỜNG HỢP 2: CHỈ MỞ CHO KHÁCH QUEN
		corsConfig.AllowOrigins = strings.Split(cfg.ClientOrigin, ",")
	}
	// Các Method được phép gọi
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	// Các Header được phép gửi lên
	//Origin: Nơi gửi đến.
	//Content-Type: Định dạng dữ liệu (JSON, XML...).
	//Authorization: Cực quan trọng. Đây là chỗ chứa Token đăng nhập. Nếu thiếu cái này, User đăng nhập xong gửi Token lên sẽ bị Server từ chối nhận.
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	return cors.New(corsConfig)
}
