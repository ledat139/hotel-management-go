package usecase

import (
	"context"
	"errors"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user *models.User) error {
	err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) UpdateUserProfile(ctx context.Context, userProfile *dto.UpdateProfileRequest, userEmail string) error {
	user, err := u.repo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return errors.New("error.get_user_failed")
	}
	if userProfile.Name != "" {
		user.Name = userProfile.Name
	}
	if userProfile.PhoneNumber != "" {
		user.PhoneNumber = userProfile.PhoneNumber
	}
	if userProfile.AvatarURL != "" {
		user.AvatarURL = userProfile.AvatarURL
	}
	err = u.repo.UpdateUser(ctx, user)
	if err != nil {
		return errors.New("error.failed_to_update_user")
	}
	return nil
}