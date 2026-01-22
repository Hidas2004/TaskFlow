package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

// 1 lưu thong tin file vào DB
func (ar *attachmentRepository) Create(attachment *models.Attachment) error {
	//db.create sẽ tự sinh câu lẹnh SQL
	return ar.db.Create(attachment).Error
}

// 2. Lấy danh sách file của một Task
func (r *attachmentRepository) FindByTaskID(taskID uuid.UUID) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	err := r.db.
		Where("task_id = ?", taskID). //điều kiện lọc
		Preload("Uploader").
		Order("created_at desc").
		Find(&attachments).Error
	return attachments, err
}

// 3. Tìm file theo ID
func (ar *attachmentRepository) FindByID(id uuid.UUID) (*models.Attachment, error) {
	var attachment models.Attachment
	err := ar.db.First(&attachment, "id = ?", id).Error
	return &attachment, err
}

func (ar *attachmentRepository) Delete(id uuid.UUID) error {
	return ar.db.Delete(&models.Attachment{}, "id = ? ", id).Error
}
