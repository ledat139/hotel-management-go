package models

import "gorm.io/gorm"

type RoomImage struct {
	gorm.Model
	RoomID   uint   `gorm:"not null" json:"room_id"`
	ImageURL string `gorm:"type:varchar(255);not null" json:"image_url" binding:"required,url"`
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}
