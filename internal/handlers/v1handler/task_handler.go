package v1handler

import (
	"errors"
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskHandler struct giữ kết nối tới Service
type TaskHandler struct {
	taskService v1services.TaskService
}

// Constructor để main.go gọi
func NewTaskHandler(service v1services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: service,
	}
}

// Hàm helper (private): Lấy UserID từ context (do Middleware nhét vào)
func (th *TaskHandler) getUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("unauthorized: user id not found in context")
	}

	// Trường hợp 1: Middleware lưu dạng UUID Object (Tốt nhất)
	if id, ok := val.(uuid.UUID); ok {
		return id, nil
	}
	if idStr, ok := val.(string); ok {
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, errors.New("invalid user id format")
		}
		return parsedID, nil
	}

	// Trường hợp 3: Kiểu dữ liệu lạ
	return uuid.Nil, errors.New("user id type invalid")
}

// Hàm helper: Lấy Role
func (th *TaskHandler) getUserRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

// POST /api/tasks
func (th *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//2 lấy ID người tạo
	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//gọi service
	resp, err := th.taskService.CreateTask(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

// GET /api/tasks?page=1&limit=10&status=todo
func (th *TaskHandler) GetTasks(c *gin.Context) {
	var req dto.TaskFilterRequest
	//1 hứng query param
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//2 lấy context user
	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := th.getUserRole(c)

	//gọi sẻvice
	tasks, total, err := th.taskService.GetTasks(c.Request.Context(), req, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//trả về kèm phân trang
	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"total": total,
			"page":  req.Page,
			"limit": req.Limit,
		},
	})
}

// PUT /api/tasks/:id
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	//1 lấy ID từ URL
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := th.getUserRole(c)

	resp, err := th.taskService.UpdateTask(c.Request.Context(), taskID, req, userID, userRole)
	if err != nil {
		// Có thể lỗi do permission hoặc not found
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// PATCH /api/tasks/:id/status
func (th *TaskHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	// Struct tạm để hứng mỗi status, không cần dùng cả DTO to bự
	var req struct {
		Status string `json:"status" binding:"required,oneof=todo in_progress review done"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := th.getUserRole(c)

	err = th.taskService.UpdateTaskStatus(c.Request.Context(), taskID, req.Status, userID, userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

// DELETE /api/tasks/:id
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := th.getUserRole(c)

	err = th.taskService.DeleteTask(c.Request.Context(), taskID, userID, userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func (th *TaskHandler) GetTaskByID(c *gin.Context) {
	// 1. Lấy ID từ URL
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}
	userID, err := th.getUserID(c)
	if err != nil {
		// Trả về 401 Unauthorized ngay, không cho chạy tiếp
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := th.getUserRole(c)

	// 3. Gọi Service
	resp, err := th.taskService.GetTaskByID(c.Request.Context(), taskID, userID, userRole)
	if err != nil {
		// Xử lý mã lỗi HTTP cho từng trường hợp (Optional nhưng nên làm)
		statusCode := http.StatusBadRequest // Mặc định
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
