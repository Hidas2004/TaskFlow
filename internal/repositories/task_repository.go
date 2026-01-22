package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/dto"
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

	// 1. Khởi tạo query
	query := r.db.Model(&models.Task{})

	// 2. Filters (Lọc dữ liệu)
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

	if keyword, ok := filters["search"]; ok && keyword != "" {
		searchTerm := "%" + keyword.(string) + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Sắp xếp: Ưu tiên Position (cho kéo thả) trước, sau đó đến ngày tạo mới nhất
	err := query.Limit(limit).
		Offset(offset).
		Order("position ASC, created_at DESC").
		Preload("Team").
		Preload("Assignee").
		Preload("Creator").
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

func (r *taskRepository) CountTasksByStatus(teamID uuid.UUID) ([]*dto.TaskCountResponse, error) {
	var results []*dto.TaskCountResponse

	// Tư duy SQL: SELECT status, count(*) as count FROM tasks WHERE team_id = ? GROUP BY status
	err := r.db.Model(&models.Task{}).
		Select("status, count(*) as count").
		Where("team_id = ?", teamID). // QUAN TRỌNG: Chỉ đếm task của team này
		Group("status").
		Scan(&results).Error

	return results, err
}
