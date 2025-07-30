package repository

import (
	"context"
	"hotel-management/internal/constant"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetAllCustomers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var staffs []models.User
	if err := r.db.WithContext(ctx).Where("role = ?", constant.STAFF).Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}
func (r *userRepository) GetAllCustomers(ctx context.Context) ([]models.User, error) {
	var customers []models.User
	if err := r.db.WithContext(ctx).Where("role = ?", constant.CUSTOMER).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
