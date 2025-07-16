package usecase

import (
	"context"
	"errors"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) Register(ctx context.Context, registerRequest *dto.RegisterRequest) (*models.User, error) {
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

	user := &models.User{
		Name:         strings.TrimSpace(registerRequest.FirstName + " " + registerRequest.LastName),
		Email:        registerRequest.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, errors.New("error.create_user_failed")
	}

	return user, nil
}

func (u *UserUseCase) Authenticate(ctx context.Context, loginRequest *dto.LoginRequest) (*models.User, error) {
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

func (u *UserUseCase) AuthenticateUserFromClaim(ctx context.Context, refreshTokenInput *dto.RefreshTokenInput) (*models.User, error) {
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
