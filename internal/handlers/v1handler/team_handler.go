package v1handler

import (
	"errors"
	"net/http"

	"github.com/Hidas2004/TaskFlow/internal/dto"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamService v1services.TeamService // Handler phụ thuộc vào Interface, không phụ thuộc code cụ thể
}

func NewTeamHandler(service v1services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: service}
}

// Hàm này giúp lấy UserID từ context một cách an toàn
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	// 1. Lấy value theo key "userID" (phải khớp y chang Middleware)
	idInterface, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("không tìm thấy user id trong context")
	}
	// 2. Ép kiểu về string trước (vì JWT lưu ID là string)
	idStr, ok := idInterface.(string)
	if !ok {
		return uuid.Nil, errors.New("kiểu dữ liệu user id không hợp lệ")
	}
	// 3. Parse từ String sang UUID
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("format user id không đúng chuẩn uuid")
	}

	return userID, nil
}

// 1. Create Team (POST /api/teams)
func (th *TeamHandler) Create(c *gin.Context) {
	// Bước 1: Hứng dữ liệu JSON gửi lên
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bước 2: Ai đang tạo? (Lấy từ Helper)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bước 3: Gọi Service
	resp, err := th.teamService.CreateTeam(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Bước 4: Trả về 201 Created
	c.JSON(http.StatusCreated, resp)
}

// 2. Get All Teams (GET /api/teams)
func (th *TeamHandler) GetAll(c *gin.Context) {
	resp, err := th.teamService.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// 3. Get Team By ID (GET /api/teams/:id)
func (th *TeamHandler) GetByID(c *gin.Context) {
	// Lấy ID từ URL
	idStr := c.Param("id")
	teamID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}

	resp, err := th.teamService.GetTeamByID(teamID)
	if err != nil {
		// Ở đây bạn có thể check kỹ hơn loại lỗi để trả 404
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// 4. Update Team (PUT /api/teams/:id)
func (th *TeamHandler) Update(c *gin.Context) {
	//lấy id can sưa
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}
	// Lấy dữ liệu cần sửa
	var req dto.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Lấy người sửa (để check quyền Leader)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	//Gọi Service & Xử lý lỗi nâng cao
	resp, err := th.teamService.UpdateTeam(teamID, &req, userID)
	if err != nil {
		// Mẹo: Check lỗi cụ thể để trả status code đúng
		if err == v1services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

}

// 5. Delete Team (DELETE /api/teams/:id)
func (th *TeamHandler) Delete(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = th.teamService.DeleteTeam(teamID, userID)
	if err != nil {
		if err == v1services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa nhóm thành công"})
}

// 6. Add Member (POST /api/teams/:id/members)
func (th *TeamHandler) AddMember(c *gin.Context) {
	//Bước 1: Xác định "Nhóm nào?"
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}
	//Bước 2: Xác định "Thêm ai?"
	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Bước 3: Xác định "Ai đang thực hiện hành động này?"
	requestUserID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	//Gọi Service xử lý nghiệp vụ
	err = th.teamService.AddMember(&req, teamID, requestUserID)
	if err != nil {
		if err == v1services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Thêm thành viên thành công"})
}

// 7. Remove Member (DELETE /api/teams/:id/members/:userId)
func (th *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}

	// Lấy ID người bị xóa từ URL
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID thành viên không hợp lệ"})
		return
	}

	requestUserID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = th.teamService.RemoveMember(teamID, targetUserID, requestUserID)
	if err != nil {
		if err == v1services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành viên thành công"})
}

// 8. Get Members (GET /api/teams/:id/members)
func (th *TeamHandler) GetMembers(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID team không hợp lệ"})
		return
	}

	members, err := th.teamService.GetMembers(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

// 9. Get My Teams (GET /api/teams/my)
func (th *TeamHandler) GetMyTeams(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	teams, err := th.teamService.GetMyTeams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}
