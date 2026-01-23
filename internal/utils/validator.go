package utils

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidateTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case models.TaskStatusTodo, models.TaskStatusInProgress, models.TaskStatusDone:
		return true
	}
	return false
}

func RegisterCustomValidators() {
	// Lấy engine validator hiện tại của Gin ra để custom
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("task_status", ValidateTaskStatus)
	}
}
