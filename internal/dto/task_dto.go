package dto

import (
	"time"

	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      string     `json:"status" binding:"omitempty,task_status"`
	TeamID      string     `json:"team_id" binding:"required,uuid"`
	AssignedTo  *string    `json:"assigned_to" binding:"omitempty,uuid"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      string `json:"status" binding:"omitempty,task_status"`

	AssignedTo *uuid.UUID `json:"assigned_to"` // Update người làm
	DueDate    *time.Time `json:"due_date"`

	// Thêm Position để hỗ trợ Kéo thả (Drag & Drop) sau này
	Position *float64 `json:"position"`
}

type TaskFilterRequest struct {
	utils.PaginationQuery
	TeamID     string `form:"team_id" binding:"omitempty,uuid"`
	Status     string `form:"status"`
	Priority   string `form:"priority"`
	AssignedTo string `form:"assigned_to" binding:"omitempty,uuid"`
	Search     string `form:"search"` // Tìm kiếm từ khóa

	Page  int `form:"page,default=1"`   // Mặc định là 1 nếu không gửi
	Limit int `form:"limit,default=10"` // Mặc định là 10
}

// 4. RESPONSE (Trả về cho Frontend)
// Chúng ta định nghĩa rõ ràng các struct con thay vì dùng interface{}

type ShortUserResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	// AvatarUrl string `json:"avatar_url"` // Nếu có
}

type ShortTeamResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Position    float64   `json:"position"`

	// Struct lồng nhau: Trả về object gọn gàng, bảo mật (không lộ password user)
	Team     ShortTeamResponse  `json:"team"`
	Assignee *ShortUserResponse `json:"assignee"` // Pointer vì có thể null
	Creator  ShortUserResponse  `json:"creator"`  // Không bao giờ null

	DueDate   *time.Time `json:"due_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type TaskCountResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type DashboardFilterRequest struct {
	TeamID string `form:"team_id" binding:"required,uuid"`
}
