package v1services

import (
	"errors"
	"mime/multipart"

	"github.com/Hidas2004/TaskFlow/internal/config"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/Hidas2004/TaskFlow/internal/utils"
	"github.com/google/uuid"
)

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	taskRepo       repositories.TaskRepository
	config         *config.Config
}

func NewAttachmentService(
	attachmentRepo repositories.AttachmentRepository,
	taskRepo repositories.TaskRepository,
	cfg *config.Config,
) AttachmentService {
	return &attachmentService{
		attachmentRepo: attachmentRepo,
		taskRepo:       taskRepo,
		config:         cfg,
	}
}

// 1. Logic Upload File
func (as *attachmentService) UploadAttachment(taskID uuid.UUID, userID uuid.UUID, file *multipart.FileHeader) (*models.Attachment, error) {
	// BƯỚC 1: Kiểm tra Task có tồn tại không?
	task, err := as.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	_ = task
	// BƯỚC 2: Validate File
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".pdf", ".docx", ".zip", ".xlsx", ".txt"}
	if err := utils.ValidateFile(file, as.config.MaxUploadSize, allowedTypes); err != nil {
		return nil, err
	}
	savePath := as.config.UploadPath + "/attachments/tasks"
	filePath, err := utils.SaveFile(file, savePath)
	if err != nil {
		return nil, errors.New("failed to save file to disk")
	}
	// BƯỚC 4: Tạo Model để chuẩn bị lưu vào DB
	attachment := &models.Attachment{
		TaskID:     taskID,
		UploadedBy: userID,
		FileName:   file.Filename,
		FilePath:   filePath, // Lưu đường dẫn tương đối: "uploads/..."
		FileSize:   file.Size,
		FileType:   file.Header.Get("Content-Type"),
	}
	// BƯỚC 5: Lưu DB
	if err := as.attachmentRepo.Create(attachment); err != nil {
		// ROLLBACK: Nếu lưu DB thất bại, XÓA ngay file vừa tạo để tránh rác
		utils.DeleteFile(filePath)
		return nil, errors.New("failed to save file info to database")
	}

	return attachment, nil

}

// 2. Logic Lấy danh sách file
func (as *attachmentService) GetAttachmentsByTaskID(taskID uuid.UUID) ([]*models.Attachment, error) {
	if _, err := as.taskRepo.FindByID(taskID); err != nil {
		return nil, errors.New("task not found")
	}
	return as.attachmentRepo.FindByTaskID(taskID)
}

// 3. Logic Xóa file
func (as *attachmentService) DeleteAttachment(attachmentID uuid.UUID, userID uuid.UUID) error {
	//b1 tìm file trong db có ko
	attachment, err := as.attachmentRepo.FindByID(attachmentID)
	if err != nil {
		return errors.New("attachment not found")
	}
	//b2 check quyền
	if attachment.UploadedBy != userID {
		return errors.New("permission denied: you are not the uploader")
	}
	//b3 xóa file
	if err := utils.DeleteFile(attachment.FilePath); err != nil {

	}
	return as.attachmentRepo.Delete(attachmentID)
}
