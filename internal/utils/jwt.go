package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 1. Định nghĩa "Ruột" của Token (Payload)
// Struct này định nghĩa những thông tin bạn muốn nhét vào trong Token
type Claims struct {
	UserID               string `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Các trường chuẩn của JWT như ngày hết hạn (exp), ngày tạo (iat)
}

// 2. Hàm Tạo Token (Generate)
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

// 3. Hàm Kiểm tra Token (Validate) - Middleware sẽ gọi hàm này
func ValidateToken(tokenString, secret string) (*Claims, error) {
	// Parse: Dịch chuỗi token mã hóa ngược lại thành struct Claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Bảo mật: Kiểm tra xem thuật toán ký có đúng là HMAC không (tránh lỗ hổng "none" algorithm)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	// Xử lý kết quả
	if err != nil {
		return nil, err
	}
	// Kiểm tra xem token có hợp lệ không (ví dụ: chưa hết hạn)
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil // Trả về thông tin user (UserID, Role...)
	}

	return nil, errors.New("invalid token")
}
