package models

import (
	"time"

	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	BookingID   uint      `gorm:"not null" json:"booking_id"`
	TotalAmount float64   `gorm:"not null" json:"total_amount"`
	ExportAt    time.Time `gorm:"type:timestamp;not null" json:"export_at" binding:"required"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}
