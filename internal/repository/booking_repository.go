package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type BookingRepository interface {
	DeleteBookingRoomByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (*bookingRepository) DeleteBookingRoomByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error {
	err := tx.WithContext(ctx).Where("room_id = ?", id).Delete(&models.BookingRoom{}).Error
	if err != nil {
		return err
	}
	return nil
}
