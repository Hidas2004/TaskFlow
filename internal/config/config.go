package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config chứa toàn bộ biến môi trường của dự án
type Config struct {
	ServerPort     string
	GinMode        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpireHours int
	UploadPath     string
	MaxUploadSize  int64
}

// LoadConfig đọc file .env và nạp vào struct Config
func LoadConfig() *Config {
	// 1. Nạp file .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("lỗi ko tim thấy file env")
	}
	// 2. Điền dữ liệu vào Struct
	return &Config{
		ServerPort:     getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""), // Quan trọng: sẽ đọc pass 123456 từ .env
		DBName:         getEnv("DB_NAME", "taskflow_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 72), // Chuyển đổi string sang int
		UploadPath:     getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize:  getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // Chuyển đổi sang int64
	}
}

func getEnv(key, fallback string) string {
	//os.lookupenv sẽ trả về 2 giá trị (value) va trạng thái true false
	//if <khởi tạo biến>; <điều kiện> { ... }
	if value, exists := os.LookupEnv(key); exists {
		return value //nếu tồn tại thì trả về giá trị
	}
	return fallback //nếu không tồn tại thì trả về giá trị mặc định
}

// getEnvAsInt lấy giá trị Int từ biến môi trường
func getEnvAsInt(key string, fallback int) int {
	//kiem tra xem biến có tồn tại không
	if valueStr, exists := os.LookupEnv(key); exists {
		// Dùng strconv để đổi từ String -> Int
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return fallback
}

// getEnvAsInt64 lấy giá trị Int64 từ biến môi trường
func getEnvAsInt64(key string, fallback int64) int64 {
	if valueStr, exists := os.LookupEnv(key); exists {
		// Dùng strconv để đổi từ String -> Int64
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return value
		}
	}
	return fallback
}
