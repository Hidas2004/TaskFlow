package v1services

import (
	"errors"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/models"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrForbidden    = errors.New("bạn không có quyền thực hiện hành động này")
	ErrTeamNotFound = errors.New("không tìm thấy nhóm")
	ErrUserNotFound = errors.New("không tìm thấy người dùng")
)

type teamService struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func NewTeamService(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// Hàm này chuyển đổi dữ liệu thô từ Database (models.Team) sang dữ liệu đẹp cho API (dto.TeamResponse).
func (ts *teamService) mapToResponse(team *models.Team) *dto.TeamResponse {
	//1 xữ lý thông tin leader
	var leaderData interface{}
	if team.Leader.ID != uuid.Nil { // Kiểm tra xem Leader có tồn tại không (UUID khác rỗng)
		leaderData = map[string]interface{}{
			"id":        team.Leader.ID,
			"full_name": team.Leader.FullName,
			"email":     team.Leader.Email,
		}
	}
	// 2. Xử lý danh sách thành viên (Members)
	var membersData []interface{}
	for _, m := range team.Members {
		if m.UserID != uuid.Nil {
			membersData = append(membersData, map[string]interface{}{
				"id":        m.User.ID,
				"full_name": m.User.FullName,
				"role":      m.Role,
				"joined_at": m.JoinedAt,
			})
		}
	}
	// 3. Xử lý con trỏ Description
	desc := ""
	if team.Description != nil { // Kiểm tra con trỏ có null không
		desc = *team.Description // Lấy giá trị của con trỏ (Dereference)
	}
	// 4. Trả về struct DTO
	return &dto.TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: desc,
		LeaderID:    team.LeaderID,
		Leader:      leaderData,
		Members:     membersData,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}
}

func (ts *teamService) CreateTeam(req *dto.CreateTeamRequest, leaderID uuid.UUID) (*dto.TeamResponse, error) {
	description := &req.Description
	//tạo dối tượng team lưu vào DB
	team := &models.Team{
		Name:        req.Name,
		Description: description,
		LeaderID:    leaderID,
	}

	var members []models.TeamMember

	members = append(members, models.TeamMember{
		UserID: leaderID,
		Role:   "leader",
	})
	// 4. Duyệt qua danh sách ID thành viên được gửi lên
	for _, memberIDStr := range req.MemberIDs {
		mID, err := uuid.Parse(memberIDStr)
		if err == nil && mID != leaderID {
			members = append(members, models.TeamMember{
				UserID: mID,
				Role:   "member",
			})
		}
	}
	team.Members = members
	//gọi repo de lưu xuong DB
	if err := ts.teamRepo.Create(team); err != nil {
		return nil, err
	}
	return ts.mapToResponse(team), nil
}

func (ts *teamService) UpdateTeam(teamID uuid.UUID, req *dto.UpdateTeamRequest, userID uuid.UUID) (*dto.TeamResponse, error) {
	//1 tìm nhóm trong DB xem có tồn tại ko
	team, err := ts.teamRepo.FindByID(teamID)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	if team.LeaderID != userID {
		return nil, ErrForbidden //// Nếu không khớp, chặn ngay
	}
	// 3. Cập nhật từng trường nếu có dữ liệu gửi lên
	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = &req.Description
	}
	//lưu thay đổi
	if err := ts.teamRepo.Update(team); err != nil {
		return nil, err
	}
	return ts.mapToResponse(team), nil
}

func (s *teamService) RemoveMember(teamID, targetUserID, requestUserID uuid.UUID) error {
	//1 lấy thông tin nhóm
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return ErrTeamNotFound
	}
	// Người yêu cầu có phải Leader không (requestUserID la người yêu cauuas)?
	isLeader := team.LeaderID == requestUserID
	// Người yêu cầu có phải đang tự xóa chính mình không?
	isSelf := targetUserID == requestUserID

	// Nếu KHÔNG phải Leader VÀ KHÔNG phải tự xóa mình -> Chặn
	if !isLeader && !isSelf {
		return ErrForbidden
	}
	// Không ai được xóa Leader ra khỏi nhóm (kể cả chính Leader cũng không được tự rời kiểu này, phải chuyển quyền trước)
	if targetUserID == team.LeaderID {
		return errors.New("không thể xóa leader khỏi nhóm")
	}
	return s.teamRepo.RemoveMember(teamID, targetUserID)
}

func (ts *teamService) GetAllTeams() ([]*dto.TeamResponse, error) {
	//lấy dl thô từ repo
	teams, err := ts.teamRepo.FindAll()
	if err != nil {
		return nil, err // Nếu lỗi database thì báo lỗi ngay
	}
	//chuẩn bị ds chứa kq trả về
	var responses []*dto.TeamResponse
	for _, team := range teams {
		responses = append(responses, ts.mapToResponse(team))
	}
	return responses, nil
}

// 3. GetTeamByID
func (s *teamService) GetTeamByID(teamID uuid.UUID) (*dto.TeamResponse, error) {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return s.mapToResponse(team), nil
}

func (ts *teamService) DeleteTeam(teamID uuid.UUID, userID uuid.UUID) error {
	team, err := ts.teamRepo.FindByID(teamID)
	if err != nil {
		return ErrTeamNotFound // Nếu nhóm không tồn tại thì báo lỗi luôn
	}
	// team.LeaderID: Chủ sở hữu thực sự của nhóm (lấy từ DB).
	// userID: Người đang yêu cầu xóa (lấy từ Token đăng nhập).
	if team.LeaderID != userID {
		return ErrForbidden
	}
	return ts.teamRepo.Delete(teamID)
}

func (ts *teamService) AddMember(req *dto.AddMemberRequest, teamID, requestUserID uuid.UUID) error {
	//kiem tra su ton tại cau nhom
	team, err := ts.teamRepo.FindByID(teamID)
	if err != nil {
		return ErrTeamNotFound
	}
	if team.LeaderID != requestUserID {
		return ErrForbidden
	}
	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("user id không hợp lệ")
	}
	//kiem tra su ton tai cua thanh vien moi
	_, err = ts.userRepo.FindByID(targetUserID)
	if err != nil {
		return ErrUserNotFound
	}
	//Thực hiện thêm vào DB
	return ts.teamRepo.AddMember(teamID, targetUserID)
}

// Hiển thị danh sách nhóm mà tôi đang tham gia
func (ts *teamService) GetMyTeams(userID uuid.UUID) ([]*dto.TeamResponse, error) {
	teams, err := ts.teamRepo.GetUserTeams(userID)
	if err != nil {
		return nil, err
	}
	var responses []*dto.TeamResponse
	for _, team := range teams {
		responses = append(responses, ts.mapToResponse(team))
	}
	return responses, nil
}

func (ts *teamService) GetMembers(teamID uuid.UUID) ([]*dto.MemberResponse, error) {
	users, err := ts.teamRepo.GetTeamMembers(teamID)
	if err != nil {
		return nil, err
	}
	//Chuyển đổi (Mapping) từ Model sang Response
	//nghĩa là lấy dữ liệu từ DB(model) -> biến thành dữ liệu gọn gàg trả về cho API(response)
	var result []*dto.MemberResponse
	for _, user := range users {
		result = append(result, &dto.MemberResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     "member", // Tạm thời để cứng, muốn chuẩn phải join bảng team_members để lấy role
		})
	}
	return result, nil
}
