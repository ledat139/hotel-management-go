package database

import (
	"hotel-management/internal/models"
	"log"
)

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.RoomImage{},
		&models.Booking{},
		&models.BookingRoom{},
		&models.Review{},
		&models.Bill{},
		&models.Shift{},
		&models.Payment{},
	)

	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
}
