package dto

import (
	"time"

	"github.com/google/uuid"
)

// 1. CreateTeamRequest: Dữ liệu Frontend gửi lên để tạo nhóm
type CreateTeamRequest struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Description string   `json:"description" binding:"max=500"`
	MemberIDs   []string `json:"member_ids" binding:"omitempty,dive,uuid"`
}

// 2. UpdateTeamRequest: Dữ liệu để cập nhật nhóm
type UpdateTeamRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// 3. AddMemberRequest: Dữ liệu khi thêm thành viên
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

// 4. TeamResponse: Dữ liệu trả về cho Frontend
type TeamResponse struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	LeaderID    uuid.UUID   `json:"leader_id"`
	Leader      interface{} `json:"leader,omitempty"`
	Members     interface{} `json:"members,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MemberResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
}
