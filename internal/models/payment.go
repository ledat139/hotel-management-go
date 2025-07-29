package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	BookingID     uint      `gorm:"not null" json:"booking_id"`
	TransactionID string    `gorm:"type:varchar(100);not null;" json:"transaction_id"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"payment_method" binding:"required"`
	PaymentStatus string    `gorm:"type:varchar(50);not null" json:"payment_status" binding:"required,oneof=success pending failed"`
	PaidAt        time.Time `gorm:"type:timestamp;not null" json:"paid_at" binding:"required"`
	TxnRef        string    `gorm:"type:varchar(100);not null" json:"txn_ref"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}
