package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID        uint      `gorm:"not null" json:"user_id"`
	BookingStatus string    `gorm:"default:'booked'" json:"booking_status" binding:"required,oneof=pending booked cancelled checked_in checked_out no_show"`
	TotalPrice    float64   `gorm:"not null" json:"total_price"`
	IsPaid        bool      `gorm:"not null" json:"is_paid"`
	StartDate     time.Time `gorm:"type:datetime;not null" json:"start_date"`
	EndDate       time.Time `gorm:"type:datetime;not null" json:"end_date"`

	BookingRooms []BookingRoom `gorm:"foreignKey:BookingID" json:"booking_rooms,omitempty"`
	Reviews      []Review      `gorm:"foreignKey:BookingID" json:"reviews,omitempty"`
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
