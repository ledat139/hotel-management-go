package admin

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidRequest       = errors.New("error.invalid_request")
	ErrInvalidPrice         = errors.New("error.invalid_price_per_night")
	ErrInvalidBedNum        = errors.New("error.invalid_bed_num")
	ErrInvalidRoomID        = errors.New("error.invalid_room_id")
	ErrInvalidMultipartForm = errors.New("error.invalid_request")
	ErrTooManyImages        = errors.New("error.too_many_images")
	ErrImageTooLarge        = errors.New("error.image_too_large")
	ErrInvalidImageType     = errors.New("error.invalid_image_type")
)

type RoomFormResult struct {
	Name        string
	Type        string
	Price       float64
	BedNum      int
	ViewType    string
	Description string
	HasAircon   bool
	IsAvailable bool
	Files       []*multipart.FileHeader
}

const (
	MaxUploadImages = 5
	MaxFileSize     = 2 * 1024 * 1024 // 2MB
)

var (
	AllowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
)

func ParseRoomForm(c *gin.Context) (*RoomFormResult, error) {
	name := strings.TrimSpace(c.PostForm("name"))
	roomType := strings.TrimSpace(c.PostForm("type"))
	priceStr := c.PostForm("price_per_night")
	bedStr := c.PostForm("bed_num")
	viewType := strings.TrimSpace(c.PostForm("view_type"))
	description := strings.TrimSpace(c.PostForm("description"))
	hasAircon := c.PostForm("has_aircon") == "on"
	isAvailable := c.PostForm("is_available") == "on"

	if name == "" || roomType == "" || priceStr == "" || bedStr == "" || viewType == "" {
		return nil, ErrInvalidRequest
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		return nil, ErrInvalidPrice
	}

	beds, err := strconv.Atoi(bedStr)
	if err != nil || beds < 1 {
		return nil, ErrInvalidBedNum
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, ErrInvalidMultipartForm
	}
	files := form.File["images"]

	// Validate image num
	if len(files) > MaxUploadImages {
		return nil, ErrTooManyImages
	}

	// Validate image files
	for _, file := range files {
		if file.Size > MaxFileSize {
			return nil, ErrImageTooLarge
		}
		if !isAllowedImageType(file) {
			return nil, ErrInvalidImageType
		}
	}

	return &RoomFormResult{
		Name:        name,
		Type:        roomType,
		Price:       price,
		BedNum:      beds,
		ViewType:    viewType,
		Description: description,
		HasAircon:   hasAircon,
		IsAvailable: isAvailable,
		Files:       files,
	}, nil
}

func isAllowedImageType(file *multipart.FileHeader) bool {
	openedFile, err := file.Open()
	if err != nil {
		return false
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		return false
	}
	contentType := http.DetectContentType(buffer)

	return AllowedImageTypes[contentType]
}
