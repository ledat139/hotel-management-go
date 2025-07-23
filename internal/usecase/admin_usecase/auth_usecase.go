package admin_usecase

import (
	"context"
	"errors"
	"hotel-management/internal/constant"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: repo}
}

func (u *AuthUseCase) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("error.get_user_failed")
	}
	if user.PasswordHash == "" {
		return nil, errors.New("error.invalid_credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("error.invalid_credentials")
	}
	if user.Role == constant.CUSTOMER {
		return nil, errors.New("error.invalid_role")
	}
	return user, nil
}
