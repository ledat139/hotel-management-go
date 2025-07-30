package usecase

import (
	"context"
	"errors"
	"hotel-management/internal/dto"

	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"math"

	"gorm.io/gorm"
)

type BookingUseCase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(bookingRepo repository.BookingRepository) *BookingUseCase {
	return &BookingUseCase{bookingRepo: bookingRepo}
}

func (u *BookingUseCase) CreateBooking(ctx context.Context, createBookingRequest *dto.CreateBookingRequest, userID uint) error {
	var bookingRooms []*models.BookingRoom
	var totalPrice float64

	db := u.bookingRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, roomID := range createBookingRequest.RoomIDs {
		if roomID <= 0 {
			return errors.New("error.invalid_room_id")
		}
		isAvailable, err := u.bookingRepo.IsAvailableRoom(ctx, tx, roomID, createBookingRequest.StartDate, createBookingRequest.EndDate)
		if err != nil || !isAvailable {
			tx.Rollback()
			return errors.New("error.room_is_not_available")
		}
		price, err := u.bookingRepo.GetPriceByRoomID(ctx, tx, roomID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("error.room_not_found")
		}
		if err != nil {
			tx.Rollback()
			return errors.New("error.failed_to_get_room_price")
		}
		bookingRooms = append(bookingRooms, &models.BookingRoom{
			RoomID: uint(roomID),
			Price:  price,
		})
		nights := int(math.Ceil(createBookingRequest.EndDate.Sub(createBookingRequest.StartDate).Hours() / 24))
		totalPrice += price * float64(nights)
	}
	booking := &models.Booking{
		UserID:        uint(userID),
		BookingStatus: "booked",
		TotalPrice:    totalPrice,
		IsPaid:        false,
		StartDate:     createBookingRequest.StartDate,
		EndDate:       createBookingRequest.EndDate,
	}
	err := u.bookingRepo.CreateBookingTx(ctx, tx, booking)
	if err != nil {
		tx.Rollback()
		return errors.New("error.failed_to_create_booking")
	}
	for _, bookingRoom := range bookingRooms {
		bookingRoom.BookingID = booking.ID
		err := u.bookingRepo.CreateBookingRoomTx(ctx, tx, bookingRoom)
		if err != nil {
			tx.Rollback()
			return errors.New("error.failed_to_create_booking")
		}
	}
	if err := tx.Commit().Error; err != nil {
		return errors.New("error.failed_to_commit_transaction")
	}
	return nil
}

func (u *BookingUseCase) GetBookingHistory(ctx context.Context, userID uint) ([]dto.BookingHistoryResponse, error) {
	var bookingHistoryResponse []dto.BookingHistoryResponse
	bookings, err := u.bookingRepo.GetBookingByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bookingHistoryResponse, errors.New("error.booking_not_found")
	}
	if err != nil {
		return bookingHistoryResponse, errors.New("error.failed_to_get_booking_history")
	}
	for _, booking := range bookings {
		var bookingRooms []dto.BookingHistoryRoom
		for _, room := range booking.BookingRooms {
			bookingRooms = append(bookingRooms, dto.BookingHistoryRoom{
				ID:     room.Room.ID,
				Name:   room.Room.Name,
				Type:   room.Room.Type,
				BedNum: room.Room.BedNum,
				Price:  room.Price,
			})
		}
		bookingHistoryResponse = append(bookingHistoryResponse, dto.BookingHistoryResponse{
			ID:         booking.ID,
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.BookingStatus,
			IsPaid:     booking.IsPaid,
			Rooms:      bookingRooms,
		})
	}
	return bookingHistoryResponse, nil
}

func (u *BookingUseCase) CancelBooking(ctx context.Context, bookingID uint, userID uint) error {
	if bookingID <= 0 {
		return errors.New("error.invalid_booking_id")
	}
	booking, err := u.bookingRepo.GetBookingByBookingIDAndUserID(ctx, bookingID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("error.booking_not_found")
	}
	if err != nil {
		return errors.New("error.failed_to_get_booking")
	}
	if booking.BookingStatus != "booked" {
		return errors.New("error.failed_to_cancel_booking")
	}
	booking.BookingStatus = "cancelled"
	if err := u.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return errors.New("error.failed_to_cancel_booking")
	}
	return nil
}
