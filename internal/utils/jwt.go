package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5" // Đảm bảo bạn đang dùng v5
)

// Struct này định nghĩa những thông tin bạn muốn nhét vào trong Token
type Claims struct {
	UserID               string `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Các trường chuẩn của JWT như ngày hết hạn (exp), ngày tạo (iat)
}

// Hàm 1: Tạo Token (Cấp vé)
func GenerateToken(userID, email, role, secret string, expireHours int) (string, error) {
	// 1. Tạo nội dung cho tấm vé (Claims)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// Thời gian hết hạn = Bây giờ + số giờ quy định
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			// Thời gian phát hành = Bây giờ
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// 2. Tạo object token với thuật toán ký là HS256 (phổ biến nhất)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Ký tên lên vé bằng "con dấu bí mật" (secret key) và trả về chuỗi token
	return token.SignedString([]byte(secret))
}

// Hàm 2: Kiểm tra Token (Soát vé)
func ValidateToken(tokenString, secret string) (*Claims, error) {
	// Parse token từ chuỗi string
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Dòng này để đảm bảo hacker không đổi thuật toán ký thành "none" để lừa hệ thống
		return []byte(secret), nil
	})

	// Nếu parse thành công và token hợp lệ
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
