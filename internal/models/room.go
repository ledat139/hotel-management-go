package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Type          string  `gorm:"type:varchar(50);not null" json:"type" binding:"required"`
	PricePerNight float64 `gorm:"not null" json:"price_per_night" binding:"required,gte=0"`
	BedNum        int     `gorm:"not null" json:"bed_num" binding:"required,gte=1"`
	HasAircon     bool    `gorm:"default:true" json:"has_aircon"`
	ViewType      string  `gorm:"type:varchar(100);not null" json:"view_type" binding:"required"`
	Description   string  `gorm:"type:text" json:"description"`
	IsAvailable   bool    `gorm:"default:true" json:"is_available"`

	Images       []RoomImage   `gorm:"foreignKey:RoomID" json:"images"`
	Reviews      []Review      `gorm:"foreignKey:RoomID" json:"reviews"`
	BookingRooms []BookingRoom `gorm:"foreignKey:RoomID" json:"booking_rooms,omitempty"`
}
