package models

import "gorm.io/gorm"

type BookingRoom struct {
	gorm.Model
	RoomID    uint    `gorm:"not null" json:"room_id"`
	BookingID uint    `gorm:"not null" json:"booking_id"`
	StartDate string  `gorm:"type:datetime;not null" json:"start_date" binding:"required"`
	EndDate   string  `gorm:"type:datetime;not null" json:"end_date" binding:"required"`
	Price     float64 `gorm:"not null" json:"price"`

	Room    Room    `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}
