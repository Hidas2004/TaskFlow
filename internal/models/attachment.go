package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"json:"id"`
	TaskID     uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	Task       Task      `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE;" json:"-"`
	UploadedBy uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	Uploader   User      `gorm:"foreignKey:UploadedBy" json:"-"`
	// json:"-" để khi API trả về file, không bị lôi theo cả thông tin User (pass, email...)
	FileName string `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath string `gorm:"type:varchar(500);not null" json:"file_path"`

	FileSize int64  `json:"file_size"`                          // Dùng int64 vì file có thể nặng > 2GB
	FileType string `gorm:"type:varchar(100)" json:"file_type"` // MIME type: image/png, application/pdf

	// 5. Timestamps
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete

}

func (a *Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
