package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;" json:"id"`
	Name        string       `gorm:"not null;type:varchar(100)" json:"name"`
	Description *string      `gorm:"type:text" json:"description"`
	LeaderID    uuid.UUID    `gorm:"not null" json:"leader_id"`
	Leader      User         `gorm:"foreignKey:LeaderID" json:"leader"`
	Members     []TeamMember `gorm:"foreignKey:TeamID" json:"members"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

// TeamMember: Bảng trung gian cho quan hệ Many-to-Many giữa Team và User
type TeamMember struct {
	ID uint `gorm:"primaryKey" json:"id"`

	TeamID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_team_member" json:"team_id"`
	UserID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_team_member" json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"user"`

	Role string `gorm:"type:varchar(20);default:'member'" json:"role"`

	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}
