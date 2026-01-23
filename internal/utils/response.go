package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Cấu trúc phản hồi chung cho toàn bộ API
type ResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`         // Thông báo dễ hiểu cho user
	Data    interface{} `json:"data,omitempty"`  // Dữ liệu (nếu có)
	Error   string      `json:"error,omitempty"` // Chi tiết lỗi (cho dev debug)
}

// SuccessResponse: Dùng khi xử lý thành công
// c: Gin Context, code: HTTP Status, message: Lời nhắn, data: Dữ liệu trả về
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, ResponseData{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse: Dùng khi xảy ra lỗi
// c: Gin Context, code: HTTP Status, message: Lời nhắn lỗi chung, err: Chi tiết lỗi (biến error hoặc string)
// 2. Hàm trả về Lỗi (Error) - Dùng khi bạn biết chính xác mã lỗi (ví dụ 400 Bad Request)
func ErrorResponse(c *gin.Context, code int, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	// Dùng Abort để ngăn các middleware khác chạy tiếp
	c.AbortWithStatusJSON(code, ResponseData{
		Success: false,
		Message: message,
		Error:   errMsg,
	})
}

func HandleServiceError(c *gin.Context, err error) {
	// Mặc định là lỗi 500 (Internal Server Error)
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"

	// Check loại lỗi
	if errors.Is(err, gorm.ErrRecordNotFound) {
		statusCode = http.StatusNotFound
		message = "Resource not found"
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		// (Ví dụ nếu sau này bạn handle lỗi trùng email)
		statusCode = http.StatusConflict
		message = "Resource already exists"
	} else {
		// Các lỗi validation logic thông thường
		// Nếu bạn muốn kỹ hơn, sau này ta sẽ định nghĩa custom error
		statusCode = http.StatusBadRequest
		message = "Bad Request"
	}

	ErrorResponse(c, statusCode, message, err)
}
