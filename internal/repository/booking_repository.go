package repository

import (
	"context"
	"hotel-management/internal/constant"
	"hotel-management/internal/models"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	DeleteBookingRoomByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error
	CreateBookingTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error
	CreateBookingRoomTx(ctx context.Context, tx *gorm.DB, bookingRoom *models.BookingRoom) error
	IsAvailableRoom(ctx context.Context, tx *gorm.DB, roomID int, startDate time.Time, endDate time.Time) (bool, error)
	GetPriceByRoomID(ctx context.Context, tx *gorm.DB, roomID int) (float64, error)
	GetBookingByUserID(ctx context.Context, userID uint) ([]models.Booking, error)
	GetBookingByBookingIDAndUserID(ctx context.Context, bookingID uint, userID uint) (*models.Booking, error)
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	GetDB() *gorm.DB
	GetBookingByID(ctx context.Context, bookingID uint) (*models.Booking, error)
	GetBookingByIDTx(ctx context.Context, tx *gorm.DB, bookingID uint) (*models.Booking, error)
	UpdateBookingTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error
	GetAllBookingsWithUser(ctx context.Context) ([]models.Booking, error)
	SearchBookings(ctx context.Context, userName, bookingStatus string) ([]models.Booking, error)
	GetActiveBookingsByRoomID(ctx context.Context, roomID int) ([]models.Booking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetDB() *gorm.DB {
	return r.db
}
func (*bookingRepository) DeleteBookingRoomByRoomIDTx(ctx context.Context, tx *gorm.DB, id int) error {
	err := tx.WithContext(ctx).Where("room_id = ?", id).Delete(&models.BookingRoom{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) CreateBookingTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	err := tx.WithContext(ctx).Create(&booking).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) CreateBookingRoomTx(ctx context.Context, tx *gorm.DB, bookingRoom *models.BookingRoom) error {
	err := tx.WithContext(ctx).Create(&bookingRoom).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) IsAvailableRoom(ctx context.Context, tx *gorm.DB, roomID int, startDate time.Time, endDate time.Time) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.BookingRoom{}).
		Where("room_id = ?", roomID).
		Joins("JOIN bookings ON bookings.id = booking_rooms.booking_id").
		Where("(? < bookings.end_date) AND (? > bookings.start_date)", startDate, endDate).
		Where("bookings.booking_status IN ?", []string{"booked", "checked_in"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *bookingRepository) GetPriceByRoomID(ctx context.Context, tx *gorm.DB, roomID int) (float64, error) {
	var room models.Room
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return 0, err
	}
	return room.PricePerNight, nil
}

func (r *bookingRepository) GetBookingByUserID(ctx context.Context, userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.WithContext(ctx).Preload("BookingRooms.Room").Where("user_id = ?", userID).Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingByBookingIDAndUserID(ctx context.Context, bookingID uint, userID uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).Where("id = ? and user_id = ?", bookingID, userID).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
func (r *bookingRepository) GetBookingByID(ctx context.Context, bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.WithContext(ctx).Preload("User").Preload("BookingRooms.Room").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}
func (r *bookingRepository) GetBookingByIDTx(ctx context.Context, tx *gorm.DB, bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := tx.WithContext(ctx).First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	err := r.db.WithContext(ctx).Updates(&booking).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) GetAllBookingsWithUser(ctx context.Context) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.WithContext(ctx).Preload("User").Find(&bookings).Error
	return bookings, err
}
func (r *bookingRepository) UpdateBookingTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	err := tx.WithContext(ctx).Updates(&booking).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) SearchBookings(ctx context.Context, userName, bookingStatus string) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.WithContext(ctx).Model(&models.Booking{}).Preload("User").Preload("BookingRooms.Room")

	if userName != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").Where("users.name LIKE ?", "%"+userName+"%")
	}

	if bookingStatus != "" {
		query = query.Where("bookings.booking_status = ?", bookingStatus)
	}

	err := query.Order("bookings.created_at DESC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetActiveBookingsByRoomID(ctx context.Context, roomID int) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.WithContext(ctx).
		Preload("User").
		Joins("JOIN booking_rooms on bookings.id = booking_rooms.booking_id").
		Where("room_id = ? AND start_date <= NOW() AND end_date >= NOW() and booking_status = ?", roomID, constant.CHECKED_IN).
		Find(&bookings).Error
	return bookings, err
}
	