package v1handler

import (
	"errors"
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	service v1services.AttachmentService
}

func NewAttachmentHandler(service v1services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

// POST /api/v1/tasks/:id/attachments
func (ah *AttachmentHandler) Upload(c *gin.Context) {
	// 1. Validate ID
	// ✅ ĐÃ SỬA: Đổi "taskId" thành "id" để khớp với route RegisterAttachmentRoutes
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Task ID", err)
		return
	}

	// 2. Validate User (Lấy từ Middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", errors.New("user ID not found in context"))
		return
	}
	userID := userIDInterface.(uuid.UUID)

	// 3. Lấy file
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File is required", err)
		return
	}

	// 4. Gọi Service
	attachment, err := ah.service.UploadAttachment(taskID, userID, file)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "File uploaded successfully", attachment)
}

// GET /api/v1/tasks/:id/attachments
func (h *AttachmentHandler) GetByTask(c *gin.Context) {
	// ✅ ĐÃ SỬA: Đổi "taskId" thành "id"
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Task ID", err)
		return
	}

	attachments, err := h.service.GetAttachmentsByTaskID(taskID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Get attachments success", attachments)
}

// DELETE /api/v1/attachments/:id
func (h *AttachmentHandler) Delete(c *gin.Context) {
	// Route delete là /attachments/:id nên chỗ này dùng "id" là đúng rồi
	idStr := c.Param("id")
	attachmentID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Attachment ID", err)
		return
	}

	userIDInterface, _ := c.Get("userID")
	userID := userIDInterface.(uuid.UUID)

	err = h.service.DeleteAttachment(attachmentID, userID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Attachment deleted successfully", nil)
}
