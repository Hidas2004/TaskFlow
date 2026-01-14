package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ TaskRepository = (*taskRepository)(nil)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) FindAll(filters map[string]interface{}, page, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// --- FIX 1: Xử lý mặc định phân trang (Tránh lỗi offset âm) ---
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Chặn không cho lấy quá nhiều làm sập server
	}

	// 1. Khởi tạo query
	query := r.db.Model(&models.Task{})

	// 2. Filters
	if teamID, ok := filters["team_id"]; ok && teamID != "" {
		query = query.Where("team_id = ?", teamID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if priority, ok := filters["priority"]; ok && priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assignedTo, ok := filters["assigned_to"]; ok && assignedTo != "" {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	// Tìm kiếm từ khóa (Title hoặc Description)
	if keyword, ok := filters["search"]; ok && keyword != "" {
		searchTerm := "%" + keyword.(string) + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// 3. Đếm tổng số lượng (Count trước khi Limit/Offset)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Preload & Phân trang
	offset := (page - 1) * limit

	// Sắp xếp: Ưu tiên Position (cho kéo thả) trước, sau đó đến ngày tạo
	err := query.Preload("Assignee").Preload("Creator").Preload("Team").
		Limit(limit).
		Offset(offset).
		Order("position ASC, created_at DESC"). // Update: Sắp xếp theo Position chuẩn Kanban
		Find(&tasks).Error

	return tasks, total, err
}

func (tr *taskRepository) Create(task *models.Task) error {
	return tr.db.Create(task).Error
}

func (tr *taskRepository) FindByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := tr.db.Preload("Assignee").
		Preload("Creator").
		Preload("Team").
		First(&task, "id = ?", id).Error
	return &task, err
}

// Cập nhật
func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

// Xóa
func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}
