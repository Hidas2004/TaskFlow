package v1services

import (
	"context"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
}
type UserService interface {
	GetAllUsers(page, limit int) ([]*models.User, int64, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error
	DeleteUser(id uuid.UUID) error
	SearchUsers(keyword string) ([]*models.User, error)
}
type TeamService interface {
	CreateTeam(req *dto.CreateTeamRequest, leaderID uuid.UUID) (*dto.TeamResponse, error)
	GetAllTeams() ([]*dto.TeamResponse, error)
	GetTeamByID(teamID uuid.UUID) (*dto.TeamResponse, error)
	UpdateTeam(teamID uuid.UUID, req *dto.UpdateTeamRequest, userID uuid.UUID) (*dto.TeamResponse, error)
	DeleteTeam(teamID uuid.UUID, userID uuid.UUID) error

	AddMember(req *dto.AddMemberRequest, teamID, requestUserID uuid.UUID) error
	RemoveMember(teamID, targetUserID, requestUserID uuid.UUID) error
	GetMyTeams(userID uuid.UUID) ([]*dto.TeamResponse, error)
	GetMembers(teamID uuid.UUID) ([]*dto.MemberResponse, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, req dto.CreateTaskRequest, creatorID uuid.UUID) (*dto.TaskResponse, error)

	GetTasks(ctx context.Context, req dto.TaskFilterRequest, userID uuid.UUID, userRole string) ([]dto.TaskResponse, int64, error)

	GetTaskByID(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, userRole string) (*dto.TaskResponse, error)
	UpdateTask(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest, userID uuid.UUID, userRole string) (*dto.TaskResponse, error)

	UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, newStatus string, userID uuid.UUID, userRole string) error

	DeleteTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, userRole string) error
}

type CommentService interface {
	CreateComment(userID uuid.UUID, taskID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetCommentsByTask(userID uuid.UUID, taskID uuid.UUID) ([]*dto.CommentResponse, error)
	UpdateComment(userID uuid.UUID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(userID uuid.UUID, commentID uuid.UUID) error
}
