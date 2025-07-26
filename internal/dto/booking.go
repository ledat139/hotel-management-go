package dto

import "time"

type SearchRoomRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	BedNum    *int      `json:"bed_num"`
	HasAircon *bool     `json:"has_aircon"`
	ViewType  *string   `json:"view_type"`
	MinPrice  *float64  `json:"min_price"`
	MaxPrice  *float64  `json:"max_price"`
}

type SearchRoomResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	PricePerNight float64  `json:"price_per_night"`
	BedNum        int      `json:"bed_num"`
	HasAircon     bool     `json:"has_aircon"`
	ViewType      string   `json:"view_type"`
	Description   string   `json:"description"`
	ImageURLs     []string `json:"image_urls"`
}
