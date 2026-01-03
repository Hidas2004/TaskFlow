package handlers

import (
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersHandler struct {
	repo repositories.UserRepository
}

func NewUsersHandler(repo repositories.UserRepository) *UsersHandler {
	return &UsersHandler{
		repo: repo,
	}
}

func (uh *UsersHandler) GetUserByUuid(ctx *gin.Context) {
	idStr := ctx.Param("id")
	//kiêm tra định dạng uuid
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
		})
		return
	}

	user, err := uh.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

func (uh *UsersHandler) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "CreateUser called",
	})
}
