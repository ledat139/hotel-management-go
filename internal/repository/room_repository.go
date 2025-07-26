package repository

import (
	"context"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type RoomRepository interface {
	FindAvailableRoom(ctx context.Context, searchRoomRequest *dto.SearchRoomRequest) ([]models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) FindAvailableRoom(ctx context.Context, searchRoomRequest *dto.SearchRoomRequest) ([]models.Room, error) {
	var rooms []models.Room
	db := r.db.WithContext(ctx).Model(&models.Room{}).
		Preload("Images").
		Where("is_available = ?", true).
		Where("id NOT IN (?)", r.db.
			Model(&models.BookingRoom{}).
			Select("room_id").
			Where("(? < end_date) AND (? > start_date)", searchRoomRequest.StartDate, searchRoomRequest.EndDate),
		)

	if searchRoomRequest.BedNum != nil {
		db = db.Where("bed_num = ?", *searchRoomRequest.BedNum)
	}
	if searchRoomRequest.HasAircon != nil {
		db = db.Where("has_aircon = ?", *searchRoomRequest.HasAircon)
	}
	if searchRoomRequest.ViewType != nil {
		db = db.Where("view_type = ?", *searchRoomRequest.ViewType)
	}
	if searchRoomRequest.MinPrice != nil && searchRoomRequest.MaxPrice != nil {
		db = db.Where("price_per_night BETWEEN ? AND ?", *searchRoomRequest.MinPrice, *searchRoomRequest.MaxPrice)
	}

	if err := db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
