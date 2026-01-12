package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dòng này đảm bảo struct teamRepository implement đủ interface
var _ TeamRepository = (*teamRepository)(nil)

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

// --- BASIC CRUD ---

func (tr *teamRepository) Create(team *models.Team) error {
	return tr.db.Create(team).Error
}

func (tr *teamRepository) FindAll() ([]*models.Team, error) {
	var teams []*models.Team
	err := tr.db.Preload("Leader").Preload("Members.User").Find(&teams).Error
	return teams, err
}

func (tr *teamRepository) FindByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := tr.db.Preload("Leader").Preload("Members.User").First(&team, "id = ?", id).Error
	return &team, err
}

func (tr *teamRepository) Update(team *models.Team) error {
	return tr.db.Save(team).Error
}

func (tr *teamRepository) Delete(id uuid.UUID) error {
	return tr.db.Delete(&models.Team{}, "id = ?", id).Error
}

// --- MEMBERS MANAGEMENT ---

func (tr *teamRepository) AddMember(teamID, userID uuid.UUID) error {
	member := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   "member",
	}
	return tr.db.Create(&member).Error
}

func (tr *teamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	return tr.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{}).Error
}

func (tr *teamRepository) GetTeamMembers(teamID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	err := tr.db.Joins("JOIN team_members ON team_members.user_id = users.id").
		Where("team_members.team_id = ?", teamID).
		Find(&users).Error
	return users, err
}

func (tr *teamRepository) GetUserTeams(userID uuid.UUID) ([]*models.Team, error) {
	var teams []*models.Team
	err := tr.db.Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Preload("Leader").
		Find(&teams).Error
	return teams, err
}
