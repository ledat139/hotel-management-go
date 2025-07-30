package admin_usecase

import (
	"context"
	"errors"
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffUseCase struct {
	userRepo repository.UserRepository
}

func NewStaffUseCase(userRepo repository.UserRepository) *StaffUseCase {
	return &StaffUseCase{userRepo: userRepo}
}

func (u *StaffUseCase) CreateStaff(ctx context.Context, req *dto.CreateStaffRequest) error {
	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		return errors.New("error.failed_to_generate_password")
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error.failed_to_hash_password")
	}
	staff := &models.User{
		Name:         req.FullName,
		Email:        req.Email,
		PhoneNumber:  req.Phone,
		Role:         constant.STAFF,
		PasswordHash: string(hashedPwd),
		IsActive:     true,
	}
	_, err = u.userRepo.CreateUser(ctx, staff)
	if err != nil {
		return errors.New("error.failed_to_create_staff")
	}
	if err := utils.SendStaffPassword(staff.Email, password); err != nil {
		return errors.New("error.failed_to_send_staff_password")
	}
	return nil
}

func (u *StaffUseCase) GetAllStaffs(ctx context.Context) ([]models.User, error) {
	return u.userRepo.GetAll(ctx)
}
func (u *StaffUseCase) GetAllCustomers(ctx context.Context) ([]models.User, error) {
	return u.userRepo.GetAllCustomers(ctx)
}

func (u *StaffUseCase) GetStaffByID(ctx context.Context, id int) (*models.User, error) {
	staff, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error.staff_not_found")
		}
		return nil, errors.New("error.failed_to_get_staff")
	}
	return staff, nil
}

func (u *StaffUseCase) UpdateStaff(ctx context.Context, req *dto.UpdateStaffRequest, staffID int) error {
	staff, err := u.userRepo.GetUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("error.staff_not_found")
		}
		return errors.New("error.failed_to_get_staff")
	}
	staff.Name = req.FullName
	staff.PhoneNumber = req.Phone
	err = u.userRepo.UpdateUser(ctx, staff)
	if err != nil {
		return errors.New("error.failed_to_update_staff")
	}
	return nil
}

func (u *StaffUseCase) DeleteStaff(ctx context.Context, id int) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return errors.New("error.Failed_to_delete_staff")
	}
	return nil
}
