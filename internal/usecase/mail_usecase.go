package usecase

import (
	"context"
	"errors"
	"hotel-management/internal/dto"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type MailUseCase struct {
	userRepo repository.UserRepository
}

func NewMailUseCase(userRepo repository.UserRepository) *MailUseCase {
	return &MailUseCase{userRepo: userRepo}
}

func (u *MailUseCase) SendVerificationEmail(ctx context.Context, mailRequest dto.MailRequest) error {
	user, err := u.userRepo.GetUserByEmail(ctx, mailRequest.Email)
	if err != nil {
		return errors.New("error.get_user_failed")
	}
	if user.IsActive {
		return errors.New("error.user_already_verified")
	}
	token, err := utils.GenerateAccessToken(user)
	if err != nil {
		return errors.New("error.generate_token_failed")
	}
	if err := utils.SendVerificationEmail(user.Email, token); err != nil {
		return errors.New("error.failed_to_send_verification_email")
	}
	return nil
}

func (u *MailUseCase) ActivateAccount(ctx context.Context, tokenString string) error {
	tokenClaims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return errors.New("error.invalid_token")
	}
	user, err := u.userRepo.GetUserByEmail(ctx, tokenClaims.Email)
	if err != nil {
		return errors.New("error.get_user_failed")
	}
	if user.IsActive {
		return errors.New("error.user_already_verified")
	}
	user.IsActive = true
	err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.New("error.failed_to_update_user")
	}
	return nil
}

func (u *MailUseCase) SendResetPassword(ctx context.Context, mailRequest dto.MailRequest) error {
	user, err := u.userRepo.GetUserByEmail(ctx, mailRequest.Email)
	if err != nil {
		return errors.New("error.get_user_failed")
	}
	resetPassword, err := utils.GenerateRandomPassword(12)
	if err != nil {
		return errors.New("error.failed_to_generate_password")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error.hash_password_failed")
	}
	user.PasswordHash = string(hashedPassword)
	err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.New("error.failed_to_update_user")
	}
	if err := utils.SendResetPassword(user.Email, resetPassword); err != nil {
		return errors.New("error.failed_to_send_verification_email")
	}
	return nil
}
