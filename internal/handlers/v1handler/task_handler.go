package v1handler

import (
	"errors"
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/Hidas2004/TaskFlow/internal/utils" // Import utils
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService v1services.TaskService
}

func NewTaskHandler(service v1services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: service}
}

// --- HELPER FUNCTIONS ---

func (th *TaskHandler) getUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("unauthorized: user id not found in context")
	}
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
	return uuid.Nil, errors.New("user id type invalid")
}

func (th *TaskHandler) getUserRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

// --- HANDLERS ---

// 1. CreateTask - Tạo Task
func (th *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	userID, err := th.getUserID(c)

	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	resp, err := th.taskService.CreateTask(c.Request.Context(), req, userID)

	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Task created successfully", resp)
}

// 2. GetTasks - Lấy danh sách (Có phân trang)
func (th *TaskHandler) GetTasks(c *gin.Context) {
	var req dto.TaskFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filter parameters", err)
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userRole := th.getUserRole(c)

	tasks, total, err := th.taskService.GetTasks(c.Request.Context(), req, userID, userRole)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	response := utils.NewPaginationResponse(tasks, req.Page, req.Limit, total)

	utils.SuccessResponse(c, http.StatusOK, "Get tasks list success", response)
}

// 3. UpdateTask
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userRole := th.getUserRole(c)

	resp, err := th.taskService.UpdateTask(c.Request.Context(), taskID, req, userID, userRole)
	if err != nil {
		utils.HandleServiceError(c, err) // Tự động xử lý lỗi 403 Forbidden nếu Service trả về
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task updated successfully", resp)
}

// 4. UpdateStatus
func (th *TaskHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,task_status"` // Nhớ dùng tag task_status ta mới tạo
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status", err)
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userRole := th.getUserRole(c)

	err = th.taskService.UpdateTaskStatus(c.Request.Context(), taskID, req.Status, userID, userRole)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task status updated successfully", nil)
}

// 5. DeleteTask
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userRole := th.getUserRole(c)

	err = th.taskService.DeleteTask(c.Request.Context(), taskID, userID, userRole)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}

// 6. GetTaskByID
func (th *TaskHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	userID, err := th.getUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userRole := th.getUserRole(c)

	resp, err := th.taskService.GetTaskByID(c.Request.Context(), taskID, userID, userRole)
	if err != nil {

		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Get task details success", resp)
}
