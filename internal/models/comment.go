package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	// 1. ID: Dùng UUID cho đồng bộ với Task và User
	ID      uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content string    `gorm:"type:text;not null" json:"content"`
	// 3. Quan hệ với Task (Comment này thuộc về Task nào?)
	TaskID uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	// Constraint: OnDelete:CASCADE -> Rất quan trọng!
	// Nghĩa là: Nếu Task bị xóa, tất cả Comment của nó cũng tự động biến mất.
	Task Task `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE;" json:"-"`
	// 4. Quan hệ với User (Ai viết comment này?)
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	// 5. Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Hỗ trợ xóa mềm (Soft Delete): Xóa rồi nhưng admin vẫn có thể khôi phục hoặc xem lại
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// tự sinh code
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
