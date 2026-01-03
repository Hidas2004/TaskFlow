package utils

import "github.com/gin-gonic/gin"

// Cấu trúc phản hồi chung cho toàn bộ API
type ResponseData struct {
	Code    int         `json:"code"`            // Mã lỗi (VD: 200, 400, 401, 500)
	Message string      `json:"message"`         // Thông báo dễ đọc cho user
	Data    interface{} `json:"data,omitempty"`  // Dữ liệu trả về (nếu thành công)
	Error   interface{} `json:"error,omitempty"` // Chi tiết lỗi (nếu thất bại)
}

// SuccessResponse: Dùng khi xử lý thành công
// c: Gin Context, code: HTTP Status, message: Lời nhắn, data: Dữ liệu trả về
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, ResponseData{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse: Dùng khi xảy ra lỗi
// c: Gin Context, code: HTTP Status, message: Lời nhắn lỗi chung, err: Chi tiết lỗi (biến error hoặc string)
func ErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	var errDetails interface{}

	// Kiểm tra nếu err là kiểu error thì lấy nội dung string, còn không thì giữ nguyên
	if e, ok := err.(error); ok {
		errDetails = e.Error()
	} else {
		errDetails = err
	}

	c.AbortWithStatusJSON(code, ResponseData{
		Code:    code,
		Message: message,
		Error:   errDetails,
	})
}
