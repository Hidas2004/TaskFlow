package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	//default:gen_random_uuid()  nếu quên ko điền ID thì nó sẽ tự tạo ID ngẫu nhiên
	//gorm:"not null"có nghĩa là "Bắt buộc phải có dữ liệu".
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Email    string    `gorm:"unique;not null;type:varchar(100)" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	FullName string    `gorm:"not null;type:varchar(100)" json:"full_name"`
	// Role mặc định là 'member'
	Role string `gorm:"type:varchar(20);default:'member';not null" json:"role"`
	// Dùng con trỏ *string để cho phép giá trị NULL
	AvatarURL *string        `gorm:"type:varchar(500)" json:"avatar_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate là một "Hook" của GORM.
// Nó sẽ tự động chạy TRƯỚC khi lệnh Create được gửi xuống Database.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	//kiem tra xem ID có rỗng ko
	if u.ID == uuid.Nil {
		//tạo mới uuid
		u.ID = uuid.New()
	}
	return
}
