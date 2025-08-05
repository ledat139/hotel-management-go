package usecase

import (
	"context"
	"hotel-management/internal/dto"
	"hotel-management/internal/repository"
)

type RoomUseCase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUseCase(roomRepo repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{roomRepo: roomRepo}
}

func (u *RoomUseCase) SearchRoom(ctx context.Context, searchRoomRequest *dto.SearchRoomRequest) ([]dto.SearchRoomResponse, error) {
	responses := make([]dto.SearchRoomResponse, 0)
	rooms, err := u.roomRepo.FindAvailableRoom(ctx, searchRoomRequest)
	if err != nil {
		return responses, err
	}

	for _, room := range rooms {
		imageURLs := make([]string, 0, len(room.Images))
		for _, img := range room.Images {
			imageURLs = append(imageURLs, img.ImageURL)
		}

		res := dto.SearchRoomResponse{
			ID:            room.ID,
			Name:          room.Name,
			Type:          room.Type,
			PricePerNight: room.PricePerNight,
			BedNum:        room.BedNum,
			HasAircon:     room.HasAircon,
			ViewType:      room.ViewType,
			Description:   room.Description,
			ImageURLs:     imageURLs,
		}
		responses = append(responses, res)
	}

	return responses, nil
}
