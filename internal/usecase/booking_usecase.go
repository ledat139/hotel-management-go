package usecase

import (
	"hotel-management/internal/repository"
)

type BookingUseCase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(bookingRepo repository.BookingRepository) *BookingUseCase {
	return &BookingUseCase{bookingRepo: bookingRepo}
}
