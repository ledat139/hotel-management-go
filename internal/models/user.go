package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null" json:"name" binding:"required,min=2"`
	Email        string `gorm:"type:varchar(150);unique;not null" json:"email" binding:"required,email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"password_hash"`
	Role         string `gorm:"type:varchar(20);not null;default:'customer'" json:"role" binding:"required,oneof=customer staff admin"`
	AvatarURL    string `gorm:"type:varchar(255)" json:"avatar_url"`
	IsActive     bool   `gorm:"default:false" json:"is_active"`
	PhoneNumber  string `gorm:"type:varchar(11);not null" json:"phone_number" binding:"required,numeric,len=10|len=11"`

	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	Payments []Payment `gorm:"foreignKey:UserID" json:"payments,omitempty"`
	Shifts   []Shift   `gorm:"foreignKey:StaffID" json:"shifts,omitempty"`
	Bills    []Bill    `gorm:"foreignKey:StaffID" json:"bills,omitempty"`
}
