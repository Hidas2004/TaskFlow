package dto

type UpdateUserRequest struct {
	FullName  string `json:"full_name" binding:"max=50"`
	AvatarURL string `json:"avatar_url"`
}
