package repositories

import (
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (ur *userRepository) FindAll(page int, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Mặc định lấy 10 item nếu không truyền limit
	}
	//Đếm tổng số bản ghi (Count)
	// Ta dùng Model(&models.User{}) để GORM biết đang đếm bảng nào
	if err := ur.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := ur.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// []*models.User
// Một danh sách (slice) các con trỏ trỏ tới User (users).
func (ur *userRepository) Search(keyword string) ([]*models.User, error) {
	var users []*models.User
	if keyword == "" {
		return []*models.User{}, nil
	}
	//% đại diện cho "bất kỳ chuỗi ký tự nào" (bao gồm cả rỗng).
	searchTerm := "%" + keyword + "%"
	err := ur.db.Where("full_name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
