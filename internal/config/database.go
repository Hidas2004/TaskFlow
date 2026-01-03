package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Hàm này nhận vào Config và trả về "Cục kết nối DB" (*gorm.DB) hoặc Lỗi (error)
func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	//1 tạo 1 chuổi thông tin kết nói dsn
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	// 2. Cấu hình Logger
	// Để khi chạy, nó hiện rõ câu lệnh SQL ra màn hình console (giúp bạn debug dễ hơn)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 3. Mở kết nối (Quan trọng!)
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Database: %w", err)
	}

	// 4. Cấu hình Connection Pool
	// Tại sao cần cái này? Để tối ưu hiệu năng, tránh việc mở quá nhiều kết nối làm sập DB.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy instance sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	fmt.Println("🚀 Kết nối PostgreSQL thành công!")
	return db, nil

}
