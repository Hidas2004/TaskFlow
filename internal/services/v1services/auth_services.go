package v1services

import (
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

// Struct service
type authService struct {
	userRepo     repositories.UserRepository
	jwtSecretKey string
}

// Constructor
func NewAuthService(repo repositories.UserRepository, secret string) AuthService {
	return &authService{
		userRepo:     repo,
		jwtSecretKey: secret,
	}
}

// Hàm Login
func (as *authService) Login(ctx *gin.Context, email, password string) (string, error) {
	// Bước 1: Tìm user
	user, err := as.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("tài khoản không tồn tại")
	}

	// Bước 2: Kiểm tra pass
	match := utils.ComparePassword(user.Password, password)
	if !match {
		return "", errors.New("sai mật khẩu")
	}

	// Bước 3: Tạo token
	// Lưu ý: user.ID trong model bạn khai báo là uuid.UUID, nên phải dùng .String() là đúng rồi
	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role, as.jwtSecretKey, 72)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Hàm Logout
func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
