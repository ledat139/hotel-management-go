package models

import "gorm.io/gorm"

type Shift struct {
	gorm.Model
	StaffID   uint   `gorm:"not null" json:"staff_id"`
	StartTime string `gorm:"type:datetime;not null" json:"start_time" binding:"required"`
	EndTime   string `gorm:"type:datetime;not null" json:"end_time" binding:"required"`

	User User `gorm:"foreignKey:StaffID" json:"user,omitempty"`
}
