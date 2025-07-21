package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	mailUseCase *usecase.MailUseCase
}

func NewMailHandler(mailUseCase *usecase.MailUseCase) *MailHandler {
	return &MailHandler{mailUseCase: mailUseCase}
}

// SendVerificationEmail godoc
// @Summary Send verification email
// @Description Send a verification email with token to user's email address if the account is not activated.
// @Tags Mail
// @Accept json
// @Produce json
// @Param request body dto.MailRequest true "Email to send verification link"
// @Success 200 {object} map[string]string "Verification email sent."
// @Failure 400 {object} map[string]string "Invalid request data. | User already verified."
// @Failure 404 {object} map[string]string "Failed to get user."
// @Failure 500 {object} map[string]string "Could not generate token. | Failed to send verification email."
// @Router /mail/smtp-verify [post]
func (h *MailHandler) SendVerificationEmail(c *gin.Context) {
	var mailRequest dto.MailRequest
	if err := c.ShouldBindJSON(&mailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	err := h.mailUseCase.SendVerificationEmail(c, mailRequest)
	if err != nil {
		switch err.Error() {
		case "error.get_user_failed":
			c.JSON(http.StatusNotFound, gin.H{"error": utils.T(c, err.Error())})
		case "error.generate_token_failed", "error.failed_to_send_verification_email":
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": utils.T(c, err.Error()),
			})
		case "error.user_already_verified":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": utils.T(c, err.Error()),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": utils.T(c, "error.internal_server_error"),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.verification_email_sent")})
}

// ActiveAccountHandler godoc
// @Summary Activate user account
// @Description Activate user account based on the token sent via email
// @Tags Mail
// @Accept json
// @Produce json
// @Param token query string true "JWT token"
// @Success 200 {object} map[string]string "Account verified successfully."
// @Failure 400 {object} map[string]string "Invalid token. | User already verified."
// @Failure 404 {object} map[string]string "Failed to get user."
// @Failure 500 {object} map[string]string "Failed to update user."
// @Router /mail/verify-account [get]
func (h *MailHandler) ActiveAccountHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.token_required")})
		return
	}
	err := h.mailUseCase.ActivateAccount(c, tokenString)
	if err != nil {
		switch err.Error() {
		case "error.get_user_failed":
			c.JSON(http.StatusNotFound, gin.H{"error": utils.T(c, err.Error())})
		case "error.invalid_token", "error.user_already_verified":
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, err.Error())})
		case "error.failed_to_update_user":
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.internal_server_error")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.account_verified_successfully")})
}

// ResetPassword godoc
// @Summary      Reset user password
// @Description  Generates a new password and sends it via email to the user.
// @Tags         Mail
// @Accept       json
// @Produce      json
// @Param        request body dto.MailRequest true "User email for password reset"
// @Success      200 {object} map[string]string "message: New password sent."
// @Failure      400 {object} map[string]string "message: Invalid request data."
// @Failure      404 {object} map[string]string "message: Failed to get user"
// @Failure      500 {object} map[string]string "error: Failed to send verification email"
// @Failure      500 {object} map[string]string "message: Failed to hash password."
// @Failure      500 {object} map[string]string "message: Failed to update user."
// @Router       /mail/reset-password [post]
func (h *MailHandler) ResetPassword(c *gin.Context) {
	var mailRequest dto.MailRequest
	if err := c.ShouldBindJSON(&mailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	err := h.mailUseCase.SendResetPassword(c, mailRequest)
	if err != nil {
		switch err.Error() {
		case "error.get_user_failed":
			c.JSON(http.StatusNotFound, gin.H{"error": utils.T(c, err.Error())})
		case "error.failed_to_generate_password", "error.hash_password_failed",
			"error.failed_to_update_user", "error.failed_to_send_verification_email":
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.internal_server_error")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.new_password_sent")})
}
