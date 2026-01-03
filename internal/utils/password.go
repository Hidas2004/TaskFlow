package utils

import "golang.org/x/crypto/bcrypt"

//mã hóa mật khẩu
//(password string): Input đầu vào là mật khẩu thô (ví dụ: "123456"), kiểu chuỗ
//string: Mật khẩu đã mã hóa (để lưu vào DB).
//error: Thông báo lỗi (nếu quá trình mã hóa bị trục trặc).
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// hàm kiểm tra compare password
//hashedPassword: Chuỗi mã hóa lấy từ Database ra (cái $2a$14$...).
//password: Mật khẩu thô người dùng nhập vào (ví dụ: "123
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
