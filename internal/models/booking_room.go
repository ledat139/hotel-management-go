package models

import (
	"gorm.io/gorm"
)

type BookingRoom struct {
	gorm.Model
	RoomID    uint    `gorm:"not null" json:"room_id"`
	BookingID uint    `gorm:"not null" json:"booking_id"`
	Price     float64 `gorm:"not null" json:"price"`

	Room    Room    `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}
