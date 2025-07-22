package dto

type UpdateProfileRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	AvatarURL   string `json:"avatar_url" binding:"omitempty,url"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,numeric,len=10|len=11"`
}
