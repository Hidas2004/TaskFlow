package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dòng này đảm bảo struct teamRepository implement đủ interface
var _ CommentRepository = (*commentRepository)(nil)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (cr *commentRepository) Create(comment *models.Comment) error {
	return cr.db.Create(comment).Error
}

func (cr *commentRepository) FindByTaskID(taskID uuid.UUID) ([]*models.Comment, error) {
	// Sửa: Dùng mảng con trỏ []*models.Comment
	var comments []*models.Comment

	err := cr.db.Where("task_id = ?", taskID).
		Preload("User").         // Kèm thông tin người post
		Order("created_at asc"). // Cũ nhất lên đầu
		Find(&comments).Error

	return comments, err
}

func (cr *commentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	//tìm commment theo id, nếu ko tìm thấy thì trả lỗi
	err := cr.db.Preload("User").First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (cr *commentRepository) Update(comment *models.Comment) error {
	return cr.db.Save(comment).Error
}

func (cr *commentRepository) Delete(id uuid.UUID) error {
	return cr.db.Delete(&models.Comment{}, "id = ?", id).Error
}
