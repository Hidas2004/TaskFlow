package v1services

import (
	"context"
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/google/uuid"
)

type taskService struct {
	taskRepo repositories.TaskRepository
	teamRepo repositories.TeamRepository
}

func NewTaskService(taskRepo repositories.TaskRepository, teamRepo repositories.TeamRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		teamRepo: teamRepo,
	}
}

// helper function(private)
// nhiệm vụ : biến model (thô) -> response (đẹp) để trả về client
func (ts *taskService) mapToResponse(t *models.Task) dto.TaskResponse {
	resp := dto.TaskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Status:    t.Status,
		Priority:  t.Priority,
		Position:  t.Position,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DueDate:   t.DueDate,
		Team: dto.ShortTeamResponse{
			ID:   t.Team.ID,
			Name: t.Team.Name,
		},
		Creator: dto.ShortUserResponse{
			ID:       t.Creator.ID,
			FullName: t.Creator.FullName,
			Email:    t.Creator.Email,
		},
	}
	//Xử lý Description (Trường có thể null)
	if t.Description != nil {
		resp.Description = *t.Description
	}
	//Xử lý Assignee (Người được giao việc - Có thể chưa có ai)
	if t.Assignee != nil {
		resp.Assignee = &dto.ShortUserResponse{
			ID:       t.Assignee.ID,
			FullName: t.Assignee.FullName,
			Email:    t.Assignee.Email,
		}
	}
	return resp
}

// Khai báo một cái Map (Bản đồ luật đi đường)
// Cấu trúc: [Nơi đang đứng] : { Những nơi được phép đi tới }
var validStatusTransitions = map[string][]string{

	// 1. Nếu đang ở ô "todo" (Mới tạo/Cần làm)
	// -> Luật: Chỉ có 1 đường duy nhất là bắt tay vào làm ("in_progress").
	// -> Không được nhảy cóc sang "review" hay "done" ngay.
	"todo": {"in_progress"},

	// 2. Nếu đang ở ô "in_progress" (Đang làm)
	// -> Luật: Có 2 ngã rẽ:
	//    - Một là quay lại "todo" (Thôi không làm nữa, trả lại kho).
	//    - Hai là làm xong rồi thì gửi đi "review" (Kiểm tra).
	"in_progress": {"todo", "review"},

	// 3. Nếu đang ở ô "review" (Đang chờ Sếp duyệt)
	// -> Luật: Có 2 kết quả:
	//    - Sếp chê -> Đuổi về "in_progress" bắt làm lại.
	//    - Sếp khen -> Cho qua cửa, sang "done" (Hoàn thành).
	"review": {"in_progress", "done"},

	// 4. Nếu đang ở ô "done" (Đã xong)
	// -> Luật: Chỉ được phép quay lại "review" (Ví dụ: phát hiện lỗi, mở lại để check).
	// -> Không được tự ý chuyển về "todo" hay "in_progress" lung tung.
	"done": {"review"},
}

// CreateTask: Tạo task mới.
func (ts *taskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest, creatorID uuid.UUID) (*dto.TaskResponse, error) {
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, errors.New("invalid team_id format")
	}
	// 2. LOGIC CHECK: Team có tồn tại không?
	_, err = ts.teamRepo.FindByID(teamID)
	if err != nil {
		return nil, errors.New("team not found")
	}
	// 3. LOGIC CHECK: Nếu có assign cho ai đó, người đó phải trong Team
	//Tại sao dùng con trỏ? Vì trong Database, cột assigned_to cho phép NULL (Task chưa giao cho ai)
	var assigneeID *uuid.UUID
	if req.AssignedTo != nil {
		parsedID, err := uuid.Parse(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("invalid assigned_to format")
		}
		// Gọi Repo để check user có thuộc team không
		isMember, err := ts.teamRepo.CheckIsMember(teamID, parsedID)
		if err != nil {
			return nil, err // Lỗi DB
		}
		if !isMember {
			return nil, errors.New("assignee is not a member of this team")
		}
		assigneeID = &parsedID
	}
	newTask := &models.Task{
		Title:       req.Title,
		Description: &req.Description,
		Priority:    req.Priority, // Hook BeforeCreate sẽ lo nếu rỗng
		TeamID:      teamID,
		CreatedBy:   creatorID,
		AssignedTo:  assigneeID,
		DueDate:     req.DueDate,
		// Status tự động là 'todo' nhờ GORM Hook
	}
	//5 save db
	if err := ts.taskRepo.Create(newTask); err != nil {
		return nil, err
	}
	fullTask, err := ts.taskRepo.FindByID(newTask.ID)
	if err != nil {
		return nil, err
	}
	resp := ts.mapToResponse(fullTask)
	return &resp, nil
}

func (ts *taskService) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, newStatus string, userID uuid.UUID, userRole string) error {
	//1 lấy task hiện tại từ DB
	task, err := ts.taskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status == newStatus {
		return nil
	}
	isValid := false
	allowedStatuses := validStatusTransitions[task.Status]
	for _, status := range allowedStatuses {
		if status == newStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid status transition from " + task.Status + " to " + newStatus)
	}
	//Chỉ có Admin hoặc Leader hoặc Assignee mới được chuyển sang 'done'
	if newStatus == "done" {
		if userRole != "admin" && task.CreatedBy != userID && (task.AssignedTo != nil && *task.AssignedTo != userID) {
			return errors.New("you imply do not have permission to mark this task as done")
		}
	}
	// 4. Update
	task.Status = newStatus
	return ts.taskRepo.Update(task)
}

// GetTasks: Lấy danh sách task (kèm bộ lọc & phân trang).
func (ts *taskService) GetTasks(ctx context.Context, req dto.TaskFilterRequest, userID uuid.UUID, userRole string) ([]dto.TaskResponse, int64, error) {
	filters := make(map[string]interface{})

	// Copy các filter từ request vào map
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	// ---------------------------------------------------------
	// 👇👇👇 BỔ SUNG ĐOẠN NÀY VÀO ĐÂY 👇👇👇
	if req.Priority != "" {
		filters["priority"] = req.Priority
	}
	if req.AssignedTo != "" {
		filters["assigned_to"] = req.AssignedTo
	}
	// ---------------------------------------------------------

	// 2. LOGIC PHÂN QUYỀN (Quan trọng)
	// Admin: Không bị ép buộc filter nào (trừ khi họ tự chọn)
	// User thường: BỊ ÉP buộc phải filter theo Team họ thuộc về
	if userRole != "admin" {
		//Lấy danh sách TeamID mà user này tham gia
		userTeams, err := ts.teamRepo.GetUserTeams(userID)
		if err != nil {
			return nil, 0, err
		}
		if len(userTeams) == 0 {
			return []dto.TaskResponse{}, 0, nil
		}
		if req.TeamID != "" {
			// Check xem user có trong team đó không
			reqTeamUUID, _ := uuid.Parse(req.TeamID)
			isMember := false
			for _, team := range userTeams {
				if team.ID == reqTeamUUID {
					isMember = true
					break
				}
			}
			if !isMember {
				return nil, 0, errors.New("permission denied: you are not in this team")
			}
			filters["team_id"] = req.TeamID
		} else {
			return nil, 0, errors.New("please specify a team_id")
		}
	} else {
		// Nếu là Admin mà có chọn TeamID thì vẫn gán vào
		if req.TeamID != "" {
			filters["team_id"] = req.TeamID
		}
	}
	// 3. Gọi Repo
	tasks, total, err := ts.taskRepo.FindAll(filters, req.Page, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	// 4. Convert sang Response DTO
	var responses []dto.TaskResponse
	for _, t := range tasks {
		responses = append(responses, ts.mapToResponse(t))
	}

	return responses, total, nil
}

func (ts *taskService) checkEditPermission(ctx context.Context, task *models.Task, userID uuid.UUID, userRole string) error {
	// 1. Admin luôn có quyền
	if userRole == "admin" {
		return nil
	}
	// 2. Creator (Người tạo task) có quyền
	if task.CreatedBy == userID {
		return nil
	}
	// 3. Leader của Team sở hữu task này có quyền
	// Lưu ý: Task model của bạn phải Preload Team rồi mới check được dòng này
	if task.Team.LeaderID == userID {
		return nil
	}
	return errors.New("permission denied: only admin, team leader, or creator can perform this action")
}

// GetTaskByID: Xem chi tiết một task.
func (ts *taskService) GetTaskByID(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, userRole string) (*dto.TaskResponse, error) {
	// 1. Tìm Task trong DB
	task, err := ts.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// 2. BẢO MẬT: Nếu không phải Admin, user phải là thành viên của Team mới được xem
	if userRole != "admin" {
		// Gọi Repo check xem user có trong team này không
		isMember, err := ts.teamRepo.CheckIsMember(task.TeamID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("permission denied: you are not a member of this team")
		}
	}

	// 3. Map sang DTO và trả về
	resp := ts.mapToResponse(task)
	return &resp, nil
}

// UpdateTask: Sửa nội dung task (Tiêu đề, Mô tả, Deadline...).
func (ts *taskService) UpdateTask(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest, userID uuid.UUID, userRole string) (*dto.TaskResponse, error) {
	// 1. Lấy Task cũ lên
	task, err := ts.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	// 2. Check quyền (Dùng hàm helper ở bước 1)
	// Chỉ Admin/Leader/Creator mới được sửa nội dung task
	if err := ts.checkEditPermission(ctx, task, userID, userRole); err != nil {
		return nil, err
	}
	// 3. Mapping dữ liệu mới (Chỉ update cái gì có gửi lên)
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		// Vì Description là pointer string, nên gán trực tiếp
		desc := req.Description
		task.Description = &desc
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	// 4. Logic AssignTask (Giao việc)
	if req.AssignedTo != nil {
		// Parse ID
		newAssigneeID := *req.AssignedTo

		// Check xem người được giao có trong team không
		isMember, err := ts.teamRepo.CheckIsMember(task.TeamID, newAssigneeID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("assignee is not in this team")
		}

		// OK thì gán
		task.AssignedTo = &newAssigneeID
	}
	if err := ts.taskRepo.Update(task); err != nil {
		return nil, err
	}
	updatedTask, _ := ts.taskRepo.FindByID(taskID)
	resp := ts.mapToResponse(updatedTask)
	return &resp, nil
}

func (ts *taskService) DeleteTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, userRole string) error {
	// 1. Lấy Task để kiểm tra quyền sở hữu
	task, err := ts.taskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	// 2. Check quyền (Dùng lại hàm helper thần thánh)
	if err := ts.checkEditPermission(ctx, task, userID, userRole); err != nil {
		return err
	}

	// 3. Gọi Repo xóa (Soft delete nhờ GORM DeletedAt)
	return ts.taskRepo.Delete(taskID)

}
