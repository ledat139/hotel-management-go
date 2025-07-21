package dto

type MailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
