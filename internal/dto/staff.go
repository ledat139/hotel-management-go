package dto

type CreateStaffRequest struct {
	FullName string `form:"full_name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Phone    string `form:"phone" binding:"required"`
}

type UpdateStaffRequest struct {
	FullName string `form:"full_name"`
	Phone    string `form:"phone"`
}
