package v1services

import (
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
)

// Struct service
type authService struct {
	userRepo repositories.UserRepository
	config   *config.Config
}

// Constructor
func NewAuthService(repo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: repo,
		config:   cfg,
	}
}

// Hàm Login
func (as *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Bước 1: Tìm user theo email
	user, err := as.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}
	// BƯỚC 2: So sánh password
	match := utils.ComparePassword(user.Password, req.Password)
	if !match {
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}
	// BƯỚC 3: Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role, as.config.JWTSecret, as.config.JWTExpireHours)
	if err != nil {
		return nil, err
	}
	// BƯỚC 4: Return AuthResponse
	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

// Hàm Logout
func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

// register
func (as *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	//buoc1 check email tồn tại
	//gọi repo tìm email
	_, err := as.userRepo.FindByEmail(req.Email)
	if err == nil {
		// Nếu KHÔNG có lỗi => Tức là tìm thấy user => Báo lỗi trùng email
		return nil, errors.New("email đã tồn tại")
	}
	//bước 2 mã hóa mk
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	//bước 3ánh xạ dto sang model
	newUser := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     "member",
	}
	//bước 4 lưu user mới vào db
	if err := as.userRepo.Create(newUser); err != nil {
		return nil, err
	}
	//bước 5 tạo token
	token, err := utils.GenerateToken(newUser.ID.String(), newUser.Email, newUser.Role, as.config.JWTSecret, as.config.JWTExpireHours)
	if err != nil {
		return nil, err
	}
	// BƯỚC 6: Return AuthResponse
	// Trả về cả Token và thông tin User cho Frontend dùng
	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{ // Map thủ công
			ID:        newUser.ID,
			Email:     newUser.Email,
			FullName:  newUser.FullName,
			Role:      newUser.Role,
			AvatarURL: newUser.AvatarURL,
		},
	}, nil

}
