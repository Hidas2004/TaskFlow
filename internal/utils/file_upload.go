package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 1 validatefile : kiểm tra xem file có hợp lệ không (kich thước , đuôi file)
func ValidateFile(file *multipart.FileHeader, maxSize int64, allowedExtensions []string) error {
	//b1 kiểm tra size
	if file.Size > maxSize {
		return errors.New("file size exceeds the limit")
	}
	//b2 kiểm tra duuoi file
	//file.Filename tên đầy đủ filepath.Ext lấy đuôi file và chuyển về chữ thường
	ext := strings.ToLower(filepath.Ext(file.Filename))
	isValid := false
	for _, allowed := range allowedExtensions {
		if ext == strings.ToLower(allowed) {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("file type not allowed")
	}
	return nil
}

// 2 savefile
func SaveFile(file *multipart.FileHeader, uploadPath string) (string, error) {
	//b1 mở file gốc
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	// defer src.Close(): Dặn Go là "khi nào chạy xong hàm này, nhớ đóng file giùm tôi".
	// Nếu quên dòng này -> RAM bị rò rỉ (Memory Leak).
	defer src.Close()
	// Bước 2: Đảm bảo thư mục lưu trữ tồn tại
	//os.MKdirAll : nếu thư mục chưa có , nó tự tạo (kể cả thư mục cha / con)
	//0755 Quyền truy cập (mình đọc/ghi được, người khác chỉ đọc).
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}
	// Bước 3: Tạo tên file mới (Unique) để không bị trùng
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	dstPath := filepath.Join(uploadPath, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return dstPath, nil

}

// 3. DeleteFile
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}
