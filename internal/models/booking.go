package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID        uint    `gorm:"not null" json:"user_id"`
	BookingStatus string  `gorm:"type:enum('booked','cancelled','checked_in','checked_out','no_show');default:'booked'" json:"booking_status"`
	TotalPrice    float64 `gorm:"not null" json:"total_price"`
	IsPaid        bool    `gorm:"not null" json:"is_paid"`

	BookingRooms []BookingRoom `gorm:"foreignKey:BookingID" json:"booking_rooms,omitempty"`
	Reviews      []Review      `gorm:"foreignKey:BookingID" json:"reviews,omitempty"`
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
