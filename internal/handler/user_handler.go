package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// HandlerUpdateProfileUser godoc
// @Summary Update user profile
// @Description Update the name, avatar, or phone number of the current user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} map[string]string "message: Update profile successful."
// @Failure 400 {object} map[string]string "message: Invalid request data."
// @Failure 404 {object} map[string]string "error: Failed to get user"
// @Failure 500 {object} map[string]string "error: Failed to update user. or Internal server error."
// @Router /users/update-profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var updateProfileRequest dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ỹà-ỹ\s]+$`)
	if !nameRegex.MatchString(updateProfileRequest.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}
	err := h.userUseCase.UpdateUserProfile(c.Request.Context(), &updateProfileRequest, userEmail.(string))
	if err != nil {
		switch err.Error() {
		case "error.get_user_failed":
			c.JSON(http.StatusNotFound, gin.H{"error": utils.T(c, err.Error())})
		case "error.failed_to_update_user":
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.internal_server")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.update_profile_successful")})
}
