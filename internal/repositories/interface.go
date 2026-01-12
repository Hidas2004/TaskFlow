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
}
