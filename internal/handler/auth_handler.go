package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userUseCase *usecase.UserUseCase
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(userUseCase *usecase.UserUseCase, authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
	}
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
	user, err := h.authUseCase.Register(c.Request.Context(), &registerRequest)
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
	user, err := h.authUseCase.Authenticate(c.Request.Context(), &loginRequest)
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
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	user, err := h.authUseCase.AuthenticateUserFromClaim(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	newAccessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.generate_token_failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// GoogleLoginHandler godoc
// @Summary Google OAuth2 Login
// @Description Redirects to Google OAuth2 login
// @Tags Auth
// @Success 307 {string} string "Temporary Redirect"
// @Router /auth/google/login [get]
func (h *AuthHandler) GoogleLoginHandler(c *gin.Context) {
	url := h.authUseCase.GetGoogleLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallbackHandler godoc
// @Summary Google OAuth2 callback
// @Description Handle Google OAuth2 callback, exchange code for tokens, fetch user info, create user if not exists, and return JWT tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} map[string]interface{} "Login successful!"
// @Failure 400 {object} map[string]string "Code not found from Google."
// @Failure 500 {object} map[string]string "Failed to exchange token from Google. / Failed to get user information from Google. / Failed to get user / Failed to create user / Could not generate token."
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.code_not_found")})
		return
	}

	userInfo, err := h.authUseCase.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(c, userInfo.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.get_user_failed")})
		return
	}
	if user == nil {
		user, err = h.userUseCase.CreateUser(c, &models.User{Email: userInfo.Email, Name: userInfo.Name, IsActive: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.create_user_failed")})
			return
		}
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
