package v1services

import (
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/google/uuid"
)

type commentService struct {
	commentRepo repositories.CommentRepository
	userRepo    repositories.UserRepository
	taskRepo    repositories.TaskRepository
	teamRepo    repositories.TeamRepository
}

func NewCommentService(
	cRepo repositories.CommentRepository,
	uRepo repositories.UserRepository,
	tRepo repositories.TaskRepository,
	teamRepo repositories.TeamRepository,
) CommentService {
	return &commentService{
		commentRepo: cRepo,
		userRepo:    uRepo,
		taskRepo:    tRepo,
		teamRepo:    teamRepo,
	}
}

// --- HÀM CHÍNH: CREATE COMMENT ---
func (cs *commentService) CreateComment(userID uuid.UUID, taskID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// B1: Check quyền (Gọi hàm helper ở dưới)
	if err := cs.checkPermission(userID, taskID); err != nil {
		return nil, err
	}

	// B2: Chuẩn bị Data Model
	newComment := models.Comment{
		Content: req.Content,
		TaskID:  taskID,
		UserID:  userID,
		// CreatedAt được GORM tự xử lý
	}

	// B3: Lưu vào DB
	if err := cs.commentRepo.Create(&newComment); err != nil {
		return nil, err
	}

	// B4: Lấy lại thông tin User để trả về cho Frontend (hiển thị Avatar, Tên)
	user, err := cs.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// B5: Map sang DTO Response
	res := &dto.CommentResponse{
		ID:        newComment.ID,
		Content:   newComment.Content,
		TaskID:    newComment.TaskID,
		UserID:    newComment.UserID,
		CreatedAt: newComment.CreatedAt,
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			AvatarURL: user.AvatarURL,
		},
	}

	return res, nil
}

// --- HÀM PHỤ: CHECK PERMISSION (Đã sửa lỗi) ---
func (cs *commentService) checkPermission(userID uuid.UUID, taskID uuid.UUID) error {
	// 1. Lấy thông tin User trước
	user, err := cs.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("không tìm thấy user")
	}

	// 2. Check quyền Admin (Quyền lực tối thượng - Return luôn nếu đúng)
	if user.Role == "admin" {
		return nil
	}

	// 3. Nếu không phải Admin, bắt đầu check logic Team
	// Lấy thông tin task để biết nó thuộc team nào
	task, err := cs.taskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task không tồn tại")
	}

	// 4. Check xem user có phải thành viên của Team chứa Task đó không
	isMember, err := cs.teamRepo.CheckIsMember(task.TeamID, userID)
	if err != nil {
		return err // Lỗi DB
	}

	if !isMember {
		return errors.New("bạn không phải thành viên team này, không thể comment")
	}

	// Nếu qua hết các cửa ải trên -> OK
	return nil
}
func (cs *commentService) GetCommentsByTask(userID uuid.UUID, taskID uuid.UUID) ([]*dto.CommentResponse, error) {
	if err := cs.checkPermission(userID, taskID); err != nil {
		return nil, err
	}
	//2 gọi repo lấy ds comment
	comment, err := cs.commentRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	var res []*dto.CommentResponse
	for _, c := range comment {
		res = append(res, &dto.CommentResponse{
			ID:      c.ID,
			Content: c.Content,
			TaskID:  c.TaskID,
			UserID:  c.UserID,
			User: dto.UserResponse{ // Map thông tin User
				ID:        c.User.ID,
				FullName:  c.User.FullName,
				AvatarURL: c.User.AvatarURL,
			},
		})
	}
	return res, nil
}

func (cs *commentService) UpdateComment(userID uuid.UUID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	//1 tìm comment xem có tồn tại hay ko
	comment, err := cs.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, errors.New("bình luận ko tồn tại")
	}
	// 2. CHECK QUYỀN CHÍNH CHỦ: Người sửa phải là người tạo
	if comment.UserID != userID {
		return nil, errors.New(" bạn ko có quyền chỉnh sữa bình luận này ")
	}
	//3 update nội dung
	comment.Content = req.Content
	if err := cs.commentRepo.Update(comment); err != nil {
		return nil, err
	}
	//4 trả về ket quả
	return &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		CreatedAt: comment.CreatedAt,
		// Lưu ý: Nếu muốn trả về full UserInfo, bạn phải preload hoặc gọi userRepo
	}, nil
}

func (cs *commentService) DeleteComment(userID uuid.UUID, commentID uuid.UUID) error {
	//1 lấy thông tin comment
	comment, err := cs.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("bình luận không tồn tại")
	}
	//2 lấy thông tin user hiên tại để xem có phải admin ko
	user, err := cs.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	// 3. CHECK QUYỀN: (Là Admin) HOẶC (Là chủ comment)
	isOwner := comment.UserID == userID
	isAdmin := user.Role == "admin"
	if !isOwner && !isAdmin {
		return errors.New("bạn không có quyền xóa bình luận này")
	}
	return cs.commentRepo.Delete(commentID)

}
