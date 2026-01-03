package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Có ý nghĩa là: "Hãy kiểm tra ngay lập tức xem struct userRepository
//
//	đã code đủ tất cả các hàm mà interface UserRepository yêu cầu
//	hay chưa."
var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}

}

func (ur *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	// First với điều kiện id trực tiếp
	if err := ur.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 1. Create: Tạo user mới
func (ur *userRepository) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepository) Update(user *models.User) error {
	return ur.db.Save(user).Error
}
func (ur *userRepository) Delete(id uuid.UUID) error {
	return ur.db.Delete(&models.User{}, "id = ?", id).Error
}

func (ur *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
