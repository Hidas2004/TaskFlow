package v1handler

import (
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

// POST /api/v1/tasks/:taskId/attachments
func (ah *AttachmentHandler) Upload(c *gin.Context) {
	//1 lấy TaskID từ URL
	taskIDStr := c.Param("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Task ID", err.Error())
		return
	}
	// 2. Lấy UserID từ Token (Middleware đã gán vào Context)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}
	userID := userIDInterface.(uuid.UUID)
	//3 lấy file từ form-date
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File is required", err.Error())
		return
	}
	attachment, err := ah.service.UploadAttachment(taskID, userID, file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Upload failed", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "File uploaded successfully", attachment)
}

// GET /api/v1/tasks/:taskId/attachments
func (h *AttachmentHandler) GetByTask(c *gin.Context) {
	taskIDStr := c.Param("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Task ID", err.Error())
		return
	}

	attachments, err := h.service.GetAttachmentsByTaskID(taskID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch attachments", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", attachments)
}

// DELETE /api/v1/attachments/:id
func (h *AttachmentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	attachmentID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Attachment ID", err.Error())
		return
	}

	// Lấy UserID để check quyền sở hữu
	userIDInterface, _ := c.Get("userID")
	userID := userIDInterface.(uuid.UUID)

	err = h.service.DeleteAttachment(attachmentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete attachment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Attachment deleted successfully", nil)
}
