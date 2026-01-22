package v1handler

import (
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/gin-gonic/gin"
)

// GET /api/dashboard/stats?team_id=...
func (th *TaskHandler) GetDashboardStats(c *gin.Context) {
	// 1. Lấy team_id từ Query Param
	var req dto.DashboardFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//lấy userid từ token
	userID, err := th.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	stats, err := th.taskService.GetTaskCounts(c.Request.Context(), req.TeamID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
