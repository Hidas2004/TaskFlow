package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware nhận vào config để lấy số lượng giới hạn động
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {

	// 1. Định nghĩa luật (Rate)
	// Dùng số từ Config, thời gian cố định là 1 phút (hoặc bạn có thể đưa time vào config nốt)
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(cfg.RateLimitRequests), // Ép kiểu về int64 vì thư viện yêu cầu
	}

	// 2. Tạo kho lưu trữ (Store)
	// Hiện tại dùng Memory (RAM). Sau này scale lớn có thể đổi thành Redis.
	store := memory.NewStore()

	// 3. Tạo instance quản lý logic
	instance := limiter.New(store, rate)

	// 4. QUAN TRỌNG: Custom thông báo lỗi (Best Practice)
	// Mặc định thư viện trả về text thô. Chúng ta cần JSON đẹp đẽ cho Frontend.
	middleware := mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  "error",
			"code":    429,
			"message": fmt.Sprintf("Bạn đã gửi quá nhiều yêu cầu. Giới hạn là %d request/phút.", cfg.RateLimitRequests),
		})
	}))

	return middleware
}
