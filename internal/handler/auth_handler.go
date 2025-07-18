package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewAuthHandler(userUseCase *usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{userUseCase: userUseCase}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "Register user successful!"
// @Failure 400 {object} map[string]string "Email already exists or invalid request data"
// @Failure 500 {object} map[string]string "Failed to hash password or create user"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	user, err := h.userUseCase.Register(c.Request.Context(), &registerRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": utils.T(c, "success.register"),
		"user":    user,
	})
}

// Login godoc
// @Summary Login
// @Description Authenticate a user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful!"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid email or password"
// @Failure 500 {object} map[string]interface{} "Could not generate token"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	user, err := h.userUseCase.Authenticate(c.Request.Context(), &loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	accessToken, errAT := utils.GenerateAccessToken(user)
	refreshToken, errRT := utils.GenerateRefreshToken(user)

	if errAT != nil || errRT != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.generate_token_failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": utils.T(c, "success.login"),
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user_id":       user.ID,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshTokenInput body dto.RefreshTokenInput true "Refresh token input"
// @Success 200 {object} map[string]interface{} "New access token generated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Failure 500 {object} map[string]string "Could not generate token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error.invalid_request"})
		return
	}
	user, err := h.userUseCase.AuthenticateUserFromClaim(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	newAccessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error.generate_token_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
