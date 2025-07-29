package admin_usecase

import (
	"context"
	"errors"
	"fmt"
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomUseCase struct {
	roomRepo    repository.RoomRepository
	bookingRepo repository.BookingRepository
	reviewRepo  repository.ReviewRepository
}

func NewRoomUseCase(roomRepo repository.RoomRepository, bookingRepo repository.BookingRepository, reviewRepo repository.ReviewRepository) *RoomUseCase {
	return &RoomUseCase{roomRepo: roomRepo, bookingRepo: bookingRepo, reviewRepo: reviewRepo}
}

func (u *RoomUseCase) saveRoomImages(ctx *gin.Context, tx *gorm.DB, roomID uint, fileHeaders []*multipart.FileHeader) ([]string, error) {
	uploadDir := constant.UploadDir
	savedFiles := []string{}
	for _, fileHeader := range fileHeaders {
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		savePath := filepath.Join(uploadDir, filename)
		savedFiles = append(savedFiles, savePath)
		if err := ctx.SaveUploadedFile(fileHeader, savePath); err != nil {
			deleteSavedFiles(savedFiles)
			return savedFiles, errors.New("error.failed_to_save_file")
		}

		roomImage := &models.RoomImage{
			RoomID:   roomID,
			ImageURL: constant.ImageURL + filename,
		}
		if err := u.roomRepo.CreateRoomImageTx(ctx.Request.Context(), tx, roomImage); err != nil {
			deleteSavedFiles(savedFiles)
			return savedFiles, errors.New("error.failed_to_save_room_image")
		}
	}
	return savedFiles, nil
}
func deleteSavedFiles(paths []string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("failed to delete file %s: %v", path, err)
		}
	}
}
func (u *RoomUseCase) CreateRoom(ctx *gin.Context, createRoomRequest *dto.CreateRoomRequest) error {
	room := &models.Room{
		Name:          createRoomRequest.Name,
		Type:          createRoomRequest.Type,
		PricePerNight: createRoomRequest.PricePerNight,
		BedNum:        createRoomRequest.BedNum,
		HasAircon:     createRoomRequest.HasAircon,
		ViewType:      createRoomRequest.ViewType,
		Description:   createRoomRequest.Description,
		IsAvailable:   createRoomRequest.IsAvailable,
	}

	db := u.roomRepo.GetDB()
	return utils.WithTransaction(db, func(tx *gorm.DB) error {
		if err := u.roomRepo.CreateRoomTx(ctx.Request.Context(), tx, room); err != nil {
			return errors.New("error.failed_to_create_room")
		}

		if len(createRoomRequest.ImageFiles) > 0 {
			savedFiles, err := u.saveRoomImages(ctx, tx, room.ID, createRoomRequest.ImageFiles)
			if err != nil {
				deleteSavedFiles(savedFiles)
				return err
			}
		}
		return nil
	})
}

func (u *RoomUseCase) UpdateRoom(ctx *gin.Context, editRoomRequest *dto.EditRoomRequest) error {
	room, err := u.roomRepo.FindRoomByID(ctx, editRoomRequest.ID)
	if err != nil {
		return errors.New("error.room_not_found")
	}
	room.Name = editRoomRequest.Name
	room.Type = editRoomRequest.Type
	room.PricePerNight = editRoomRequest.PricePerNight
	room.BedNum = editRoomRequest.BedNum
	room.HasAircon = editRoomRequest.HasAircon
	room.ViewType = editRoomRequest.ViewType
	room.Description = editRoomRequest.Description
	room.IsAvailable = editRoomRequest.IsAvailable

	db := u.roomRepo.GetDB()
	return utils.WithTransaction(db, func(tx *gorm.DB) error {
		//Update room information
		if err := u.roomRepo.UpdateRoomTx(ctx.Request.Context(), tx, room); err != nil {
			return errors.New("error.failed_to_update_room")
		}
		//Delete room images
		for _, imageID := range editRoomRequest.ImageDeletes {
			image, err := u.roomRepo.FindRoomImageByID(ctx, imageID)
			if err != nil {
				return errors.New("error.room_image_not_found")
			}
			err = u.roomRepo.DeleteRoomImageTx(ctx, tx, imageID)
			if err != nil {
				return errors.New("error.failed_to_delete_image")
			}
			filename := filepath.Base(image.ImageURL)
			_ = os.Remove(filepath.Join(constant.UploadDir, filename))
		}
		//Upload new images
		if len(editRoomRequest.ImageFiles) > 0 {
			savedFiles, err := u.saveRoomImages(ctx, tx, room.ID, editRoomRequest.ImageFiles)
			if err != nil {
				deleteSavedFiles(savedFiles)
				return err
			}
		}
		return nil
	})
}

func (u *RoomUseCase) GetAllRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := u.roomRepo.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
func (u *RoomUseCase) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
	room, err := u.roomRepo.FindRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (u *RoomUseCase) DeleteRoom(ctx context.Context, id int) error {
	db := u.roomRepo.GetDB()
	return utils.WithTransaction(db, func(tx *gorm.DB) error {
		//1. Delete Room Images
		images, err := u.roomRepo.FindRoomImageByRoomID(ctx, id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("error.failed_to_get_room_images")
		}
		err = u.roomRepo.DeleteRoomImagesByRoomIDTx(ctx, tx, id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("error.failed_to_delete_images")
		}
		for _, img := range images {
			filename := filepath.Base(img.ImageURL)
			_ = os.Remove(filepath.Join(constant.UploadDir, filename))
		}

		//2. Delete Reviews
		if err := u.reviewRepo.DeleteByRoomIDTx(ctx, tx, id); err != nil {
			return errors.New("error.failed_to_delete_review")
		}

		//3. Delete BookingRooms
		if err := u.bookingRepo.DeleteBookingRoomByRoomIDTx(ctx, tx, id); err != nil {
			return errors.New("error.failed_to_delete_booking")
		}
		//4. Delete Room
		if err := u.roomRepo.DeleteRoomTx(ctx, tx, id); err != nil {
			return errors.New("error.failed_to_delete_room")
		}

		return nil
	})
}
