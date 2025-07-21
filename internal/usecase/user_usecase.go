package usecase

import (
	"context"
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
