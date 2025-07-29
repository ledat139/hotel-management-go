package dto

import "mime/multipart"

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
