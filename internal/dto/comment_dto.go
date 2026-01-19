package dto

import (
	"time"

	"github.com/google/uuid"
)

// 1 người dùng chỉ gửi nội dung
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CommentResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	TaskID    uuid.UUID `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
	// Nhúng thông tin User vào để hiển thị tên/avatar người chat
	User UserResponse `json:"user"`
}
