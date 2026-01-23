package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TaskStatusTodo       = "todo"        //Cần làm / Chưa bắt đầu.
	TaskStatusInProgress = "in_progress" //Đang thực hiện / Đang làm.
	TaskStatusDone       = "done"
)

// Định nghĩa luôn cho Priority để sau này dùng
const ( //(Độ ưu tiên)
	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
	TaskPriorityUrgent = "urgent"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description"` // Cho phép null

	// Trạng thái: todo, in_progress, done
	Status string `gorm:"type:varchar(50);default:'todo';index" json:"status"`
	// Mức độ ưu tiên: low, medium, high, urgent
	Priority string `gorm:"type:varchar(50);default:'medium'" json:"priority"`

	// Thứ tự hiển thị trên bảng (Dùng cho Drag & Drop sau này)
	Position float64 `gorm:"type:double precision;default:0" json:"position"`

	// --- RELATIONSHIPS ---

	// 1. Thuộc về Team nào (Bắt buộc)
	// Thêm index để query lấy task của team cho nhanh
	TeamID uuid.UUID `gorm:"type:uuid;not null;index" json:"team_id"`
	Team   Team      `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE;" json:"-"`
	// OnDelete:CASCADE -> Xóa Team là Task bay màu theo (đỡ rác DB)

	// 2. Giao cho ai (Có thể chưa giao - Nullable)
	AssignedTo *uuid.UUID `gorm:"type:uuid;index" json:"assigned_to"`
	Assignee   *User      `gorm:"foreignKey:AssignedTo" json:"assignee"`

	// 3. Ai tạo (Bắt buộc)
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Creator   User      `gorm:"foreignKey:CreatedBy" json:"creator"`

	// --- TIME ---
	DueDate *time.Time `json:"due_date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Hook: Tự động điền dữ liệu mặc định trước khi tạo
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.Status == "" {
		t.Status = TaskStatusTodo // Dùng hằng số
	}
	if t.Priority == "" {
		t.Priority = TaskPriorityMedium // Dùng hằng số
	}
	return
}
