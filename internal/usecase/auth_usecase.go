package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"
	"net/http"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	repo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (u *AuthUseCase) Register(ctx context.Context, registerRequest *dto.RegisterRequest) (*models.User, error) {
	userEmail, err := u.repo.GetUserByEmail(ctx, registerRequest.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("error.internal_server")
	}
	if userEmail != nil {
		return nil, errors.New("error.email_exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error.hash_password_failed")
	}
	userEmailInput := strings.ToLower(strings.TrimSpace(registerRequest.Email))
	user := &models.User{
		Name:         strings.TrimSpace(registerRequest.FirstName + " " + registerRequest.LastName),
		Email:        userEmailInput,
		PasswordHash: string(hashedPassword),
	}

	if _, err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, errors.New("error.create_user_failed")
	}

	return user, nil
}

func (u *AuthUseCase) Authenticate(ctx context.Context, loginRequest *dto.LoginRequest) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, errors.New("error.invalid_credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		return nil, errors.New("error.invalid_credentials")
	}
	return user, nil
}

func (u *AuthUseCase) AuthenticateUserFromClaim(ctx context.Context, refreshTokenInput *dto.RefreshTokenInput) (*models.User, error) {
	claims, err := utils.ValidateToken(refreshTokenInput.RefreshToken)
	if err != nil {
		return nil, errors.New("error.expired_or_invalid_refresh_token")
	}
	user, err := u.repo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, errors.New("error.invalid_user_refresh_token")
	}
	return user, nil
}

func (u *AuthUseCase) GetGoogleLoginURL() string {
	return utils.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (u *AuthUseCase) HandleGoogleCallback(code string) (*dto.GoogleUserInfo, error) {
	token, err := utils.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("error.failed_to_exchange_token")
	}

	client := utils.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("error.failed_to_get_user_info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error.failed_to_get_user_info")
	}

	var userInfo dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.New("error.failed_to_get_user_info")
	}

	return &userInfo, nil
}
