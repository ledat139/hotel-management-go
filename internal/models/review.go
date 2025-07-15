package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID    uint   `gorm:"not null" json:"user_id"`
	BookingID uint   `gorm:"not null" json:"booking_id"`
	RoomID    uint   `gorm:"not null" json:"room_id"`
	Rating    int    `gorm:"not null" json:"rating" binding:"required,min=1,max=5"`
	Comment   string `gorm:"type:text" json:"comment"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Room    Room    `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}
