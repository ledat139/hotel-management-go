package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type BillRepository interface {
	CreateBillTx(ctx context.Context, tx *gorm.DB, bill *models.Bill) error
	SearchBills(ctx context.Context, userName string, bookingID int, exportDate string) ([]models.Bill, error)
}
type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) CreateBillTx(ctx context.Context, tx *gorm.DB, bill *models.Bill) error {
	return tx.WithContext(ctx).Create(&bill).Error
}

func (r *billRepository) SearchBills(ctx context.Context, userName string, bookingID int, exportDate string) ([]models.Bill, error) {
	var bills []models.Bill
	query := r.db.WithContext(ctx).Model(&models.Bill{}).
		Joins("JOIN bookings ON bills.booking_id = bookings.id").
		Joins("JOIN users ON bookings.user_id = users.id")
	if userName != "" {
		query = query.Where("users.name LIKE ?", "%"+userName+"%")
	}
	if bookingID != 0 {
		query = query.Where("bills.booking_id = ?", bookingID)
	}
	if exportDate != "" {
		query = query.Where("DATE(bills.export_at) = ?", exportDate)
	}
	err := query.Preload("Booking.User").Find(&bills).Error
	return bills, err
}
