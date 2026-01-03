package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	// binding:"required": Nếu thiếu trường này -> Lỗi ngay lập tức
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// user response dùng để trả về thông tin user ra ngoài
// nó sạch hơn models.user vì ko có password và Gorm
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url"`
}

//AuthResponse dùng để trả về thông tin khi đăng nhập hoặc đăng ký
type AuthResponse struct {
	Token string       `json:"token"` // JWT Token để user dùng cho các request sau
	User  UserResponse `json:"user"`  // Thông tin user (đã được lọc sạch)
}
