package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	BookingID     uint   `gorm:"not null" json:"booking_id"`
	UserID        uint   `gorm:"not null" json:"user_id"`
	PaymentMethod string `gorm:"type:varchar(50);not null" json:"payment_method" binding:"required"`
	PaymentStatus string `gorm:"type:varchar(50);not null" json:"payment_status" binding:"required,oneof=success pending failed"`
	PaidAt        string `gorm:"type:timestamp;not null" json:"paid_at" binding:"required"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
