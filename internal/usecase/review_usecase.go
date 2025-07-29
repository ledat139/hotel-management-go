package usecase

import (
	"context"
	"errors"
	"hotel-management/internal/dto"
	reviewError "hotel-management/internal/error"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"

	"gorm.io/gorm"
)

type ReviewUseCase struct {
	bookingRepo repository.BookingRepository
	reviewRepo  repository.ReviewRepository
}

func NewReviewUseCase(bookingRepo repository.BookingRepository, reviewRepo repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{
		bookingRepo: bookingRepo,
		reviewRepo:  reviewRepo,
	}
}

func (u *ReviewUseCase) CreateReview(ctx context.Context, createReviewRequest *dto.CreateReviewRequest, userID uint) error {
	booking, err := u.bookingRepo.GetBookingByID(ctx, createReviewRequest.BookingID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return reviewError.ErrBookingNotFound
	}
	if err != nil || booking.UserID != userID {
		return reviewError.ErrFailedToGetBooking
	}
	if booking.BookingStatus != "checked_out" {
		return reviewError.ErrBookingNotCheckedOut
	}
	exists, err := u.reviewRepo.ExistsByBookingID(ctx, createReviewRequest.BookingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return reviewError.ErrReviewCheckFailed
	}
	if exists {
		return reviewError.ErrReviewAlreadyExists
	}
	review := &models.Review{
		UserID:    userID,
		BookingID: createReviewRequest.BookingID,
		RoomID:    createReviewRequest.RoomID,
		Rating:    createReviewRequest.Rating,
		Comment:   createReviewRequest.Comment,
	}
	if err := u.reviewRepo.CreateReview(ctx, review); err != nil {
		return reviewError.ErrFailedToCreateReview
	}
	return nil
}
