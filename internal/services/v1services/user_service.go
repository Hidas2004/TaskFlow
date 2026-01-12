package v1services

import (
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/google/uuid"
)

// 1. Struct giữ kết nối với Repository
type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// --- LOGIC ĐỌC DỮ LIỆU (Gọi Repo là chính) ---
func (s *userService) GetAllUsers(page, limit int) ([]*models.User, int64, error) {
	return s.userRepo.FindAll(page, limit)
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) SearchUsers(keyword string) ([]*models.User, error) {
	return s.userRepo.Search(keyword)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error {
	// Bước 1: Lấy user cũ từ DB lên
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.AvatarURL != "" {
		avt := req.AvatarURL
		user.AvatarURL = &avt
	}
	return s.userRepo.Update(user)
}
