package handler

import (
	"hotel-management/internal/usecase"
)

type BookingHandler struct {
	bookingUseCase *usecase.BookingUseCase
}

func NewBookingHandler(bookingUseCase *usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUseCase: bookingUseCase}
}
