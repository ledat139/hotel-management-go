package dto

import (
	"hotel-management/internal/models"
	"mime/multipart"
)

type CreateRoomRequest struct {
	Name          string
	Type          string
	PricePerNight float64
	BedNum        int
	HasAircon     bool
	ViewType      string
	Description   string
	IsAvailable   bool
	ImageFiles    []*multipart.FileHeader
}

type EditRoomRequest struct {
	ID            int
	Name          string
	Type          string
	PricePerNight float64
	BedNum        int
	HasAircon     bool
	ViewType      string
	Description   string
	IsAvailable   bool
	ImageFiles    []*multipart.FileHeader
	ImageDeletes  []int
}

type RoomQuery struct {
	Name      string  `form:"name"`
	HasAircon string  `form:"has_aircon"`
	MinPrice  float64 `form:"min_price"`
	MaxPrice  float64 `form:"max_price"`
}

type RoomDetailResponse struct {
	Room           *models.Room
	ActiveBookings []models.Booking
}
