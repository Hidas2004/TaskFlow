package v1services

import (
	"context"
	"testing"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =========================================================================================
// PHẦN 1: MOCK OBJECTS (DIỄN VIÊN ĐÓNG THẾ)
// =========================================================================================

// --- 1. MockTeamRepo ---
type MockTeamRepo struct {
	mock.Mock
}

// Các hàm CÓ DÙNG trong test (Viết logic mock)
func (m *MockTeamRepo) FindByID(id uuid.UUID) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepo) CheckIsMember(teamID, userID uuid.UUID) (bool, error) {
	args := m.Called(teamID, userID)
	return args.Bool(0), args.Error(1)
}

// Các hàm KHÔNG DÙNG (Return rỗng để trình biên dịch không báo lỗi)
func (m *MockTeamRepo) Create(team *models.Team) error                          { return nil }
func (m *MockTeamRepo) FindAll() ([]*models.Team, error)                        { return nil, nil }
func (m *MockTeamRepo) Update(team *models.Team) error                          { return nil }
func (m *MockTeamRepo) Delete(id uuid.UUID) error                               { return nil }
func (m *MockTeamRepo) AddMember(teamID, userID uuid.UUID) error                { return nil }
func (m *MockTeamRepo) RemoveMember(teamID, userID uuid.UUID) error             { return nil }
func (m *MockTeamRepo) GetTeamMembers(teamID uuid.UUID) ([]*models.User, error) { return nil, nil }
func (m *MockTeamRepo) GetUserTeams(userID uuid.UUID) ([]*models.Team, error)   { return nil, nil }

// --- 2. MockTaskRepo ---
type MockTaskRepo struct {
	mock.Mock
}

// Các hàm CÓ DÙNG trong test
func (m *MockTaskRepo) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepo) FindByID(id uuid.UUID) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

// Các hàm KHÔNG DÙNG (Return rỗng để trình biên dịch không báo lỗi)
func (m *MockTaskRepo) Update(task *models.Task) error { return nil }
func (m *MockTaskRepo) Delete(id uuid.UUID) error      { return nil }
func (m *MockTaskRepo) FindAll(filters map[string]interface{}, page, limit int) ([]*models.Task, int64, error) {
	return nil, 0, nil
}
func (m *MockTaskRepo) CountTasksByStatus(teamID uuid.UUID) ([]*dto.TaskCountResponse, error) {
	return nil, nil
}

// =========================================================================================
// PHẦN 2: TEST CASES (KỊCH BẢN)
// =========================================================================================

func TestCreateTask_Success(t *testing.T) {
	// --- ARRANGE (CHUẨN BỊ) ---
	mockTaskRepo := new(MockTaskRepo)
	mockTeamRepo := new(MockTeamRepo)
	service := NewTaskService(mockTaskRepo, mockTeamRepo)

	teamID := uuid.New()
	creatorID := uuid.New()
	mockTeam := &models.Team{ID: teamID, Name: "Golang Team"}

	req := dto.CreateTaskRequest{
		Title:       "Học Unit Test",
		Description: "Khó nhưng mà vui",
		Priority:    "High",
		TeamID:      teamID.String(),
	}

	// Mock hành vi (Stubbing)
	mockTeamRepo.On("FindByID", teamID).Return(mockTeam, nil)
	mockTaskRepo.On("Create", mock.AnythingOfType("*models.Task")).Return(nil)

	// Mock bước lấy lại task sau khi tạo để trả về response
	mockTaskRepo.On("FindByID", mock.Anything).Return(&models.Task{
		ID:      uuid.New(),
		Title:   req.Title,
		Team:    *mockTeam,
		Creator: models.User{ID: creatorID, FullName: "Hung Nguyen"},
	}, nil)

	// --- ACT (DIỄN) ---
	resp, err := service.CreateTask(context.Background(), req, creatorID)

	// --- ASSERT (KIỂM TRA) ---
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, "Golang Team", resp.Team.Name)

	mockTeamRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

// Thử luôn trường hợp có Assignee nhé
func TestCreateTask_WithAssignee_Success(t *testing.T) {
	// 1. Arrange
	mockTaskRepo := new(MockTaskRepo)
	mockTeamRepo := new(MockTeamRepo)
	service := NewTaskService(mockTaskRepo, mockTeamRepo)

	teamID := uuid.New()
	assigneeID := uuid.New()
	assigneeIDStr := assigneeID.String()

	// --- SỬA LỖI Ở ĐÂY: Khai báo mockTeam ---
	mockTeam := &models.Team{ID: teamID, Name: "Test Team"}

	req := dto.CreateTaskRequest{
		Title:      "Task Giao Cho Hùng",
		TeamID:     teamID.String(),
		AssignedTo: &assigneeIDStr,
	}

	// 2. Mocking
	// Trả về team giả đã tạo
	mockTeamRepo.On("FindByID", teamID).Return(mockTeam, nil)

	// Mock CheckIsMember trả về true
	mockTeamRepo.On("CheckIsMember", teamID, assigneeID).Return(true, nil)

	mockTaskRepo.On("Create", mock.Anything).Return(nil)

	// Mock FindByID trả về task đầy đủ (Kèm cả User Assignee và Team)
	mockTaskRepo.On("FindByID", mock.Anything).Return(&models.Task{
		ID:         uuid.New(),
		Title:      req.Title,
		AssignedTo: &assigneeID,

		// QUAN TRỌNG: Fake dữ liệu User để không bị lỗi nil pointer dereference
		Assignee: &models.User{
			ID:       assigneeID,
			FullName: "Hùng Developer",
			Email:    "hung@example.com",
		},

		// QUAN TRỌNG: Gán Team vào để tránh lỗi mockTeam undefined
		Team: *mockTeam,
	}, nil)

	// 3. Act
	resp, err := service.CreateTask(context.Background(), req, uuid.New())

	// 4. Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Xóa dấu * ở trước resp.Assignee.ID
	assert.Equal(t, assigneeID, resp.Assignee.ID)

	mockTeamRepo.AssertExpectations(t)
}
