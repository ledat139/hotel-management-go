package admin_usecase

import (
	"context"
	"errors"
	"hotel-management/internal/constant"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"

	"gorm.io/gorm"
)

type BookingUseCase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(bookingRepo repository.BookingRepository) *BookingUseCase {
	return &BookingUseCase{bookingRepo: bookingRepo}
}

func (u *BookingUseCase) GetAllBookingsWithUser(ctx context.Context) ([]models.Booking, error) {
	return u.bookingRepo.GetAllBookingsWithUser(ctx)
}

func (u *BookingUseCase) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	return u.bookingRepo.UpdateBooking(ctx, booking)
}

func (u *BookingUseCase) GetBookingDetail(ctx context.Context, id uint) (*models.Booking, error) {
	return u.bookingRepo.GetBookingByID(ctx, id)
}

func (u *BookingUseCase) UpdateBookingStatus(ctx context.Context, bookingID uint, status string) error {
	booking, err := u.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("error.booking_not_found")
		}
		return errors.New("error.failed_to_get_booking")
	}
	if booking.BookingStatus == status {
		return nil
	}
	if booking.BookingStatus == constant.CHECKED_OUT {
		return errors.New("error.failed_to_update_booking_status_because_checked_out")
	}
	if booking.BookingStatus == constant.CANCELLED {
		return errors.New("error.failed_to_update_booking_status_because_cancelled")
	}
	if booking.BookingStatus == constant.NO_SHOW {
		return errors.New("error.failed_to_update_booking_status_because_no_show")
	}
	booking.BookingStatus = status
	err = u.UpdateBooking(ctx, booking)
	if err != nil {
		return errors.New("error.failed_to_update_booking")
	}
	return nil
}

func (u *BookingUseCase) SearchBookings(ctx context.Context, userName, bookingStatus string) ([]models.Booking, error) {
	var bookings []models.Booking
	bookings, err := u.bookingRepo.SearchBookings(ctx, userName, bookingStatus)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error.booking_not_found")
		}
		return nil, errors.New("error.failed_to_get_booking")
	}
	return bookings, nil
}
