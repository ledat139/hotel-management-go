package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *models.Review) error
	ExistsByBookingID(ctx context.Context, bookingID uint) (bool, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) CreateReview(ctx context.Context, review *models.Review) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return err
	}
	return nil
}

func (r *reviewRepository) ExistsByBookingID(ctx context.Context, bookingID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("booking_id = ?", bookingID).
		Count(&count).Error
	return count > 0, err
}
