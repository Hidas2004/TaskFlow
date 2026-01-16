package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	FindAll(page int, limit int) ([]*models.User, int64, error)
	Search(keyword string) ([]*models.User, error)
}

type TeamRepository interface {
	Create(team *models.Team) error
	FindAll() ([]*models.Team, error)
	FindByID(id uuid.UUID) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uuid.UUID) error
	AddMember(teamID, userID uuid.UUID) error
	RemoveMember(teamID, userID uuid.UUID) error
	GetTeamMembers(teamID uuid.UUID) ([]*models.User, error)
	GetUserTeams(userID uuid.UUID) ([]*models.Team, error)
	CheckIsMember(teamID, userID uuid.UUID) (bool, error)
}

type TaskRepository interface {
	Create(task *models.Task) error

	// Hàm này cân tất cả: Lọc theo Team, Status, Assignee, Search từ khóa, Phân trang
	FindAll(filters map[string]interface{}, page, limit int) ([]*models.Task, int64, error)

	FindByID(id uuid.UUID) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id uuid.UUID) error
}

type CommentRepository interface {
	Create(comment *models.Comment) error

	// Lấy comment theo Task, nhưng phải kèm thông tin User
	FindByTaskID(taskID uuid.UUID) ([]*models.Comment, error)

	FindByID(id uuid.UUID) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uuid.UUID) error
}
