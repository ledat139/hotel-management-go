package models

import "gorm.io/gorm"

type Bill struct {
	gorm.Model
	BookingID   uint    `gorm:"not null" json:"booking_id"`
	StaffID     uint    `gorm:"not null" json:"staff_id"`
	TotalAmount float64 `gorm:"not null" json:"total_amount"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Staff   User    `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
}
