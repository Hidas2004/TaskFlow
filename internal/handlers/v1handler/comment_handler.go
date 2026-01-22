package v1handler

import (
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentService v1services.CommentService
}

func NewCommentHandler(service v1services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: service}
}

// @Route POST /api/v1/tasks/:taskId/comments
func (ch *CommentHandler) CreateComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID không hợp lệ"})
		return
	}
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ch.commentService.CreateComment(userID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

// @Route GET /api/v1/tasks/:taskId/comments
func (ch *CommentHandler) GetComments(c *gin.Context) {
	taskIDStr := c.Param("id")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID không hợp lệ"})
		return
	}

	// Kiểm tra user có tồn tại không trước khi ép kiểu để an toàn hơn
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	res, err := ch.commentService.GetCommentsByTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// @Route PUT /api/v1/comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	idStr := c.Param("id") //id của comment
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID lỗi"})
		return
	}
	//Xác định "Ai là người đang sửa?" (Authentication)
	//c.MustGet("userID") lấy ID người dùng đã được giải mã từ token JWT
	userID := c.MustGet("userID").(uuid.UUID)
	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.commentService.UpdateComment(userID, commentID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công", "data": res})
}

// @Route DELETE /api/v1/comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID lỗi"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.commentService.DeleteComment(userID, commentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
