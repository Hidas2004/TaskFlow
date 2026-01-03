package v1services

import "github.com/gin-gonic/gin"

type AuthService interface {
	// Sửa dòng này: Thêm string vào kết quả trả về
	Login(ctx *gin.Context, email, password string) (string, error)

	Logout(ctx *gin.Context) error
}
