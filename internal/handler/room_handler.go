package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUseCase *usecase.RoomUseCase
}

func NewRoomHandler(roomUseCase *usecase.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUseCase: roomUseCase}
}

// FindAvailableRoom godoc
// @Summary      Search available rooms
// @Description  Find all available rooms that match the search criteria and are not booked during the requested time range
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param        request body dto.SearchRoomRequest true "Search filters for room availability"
// @Success      200 {object} map[string][]dto.SearchRoomResponse "Find available room successful!"
// @Failure      400 {object} map[string]string "Invalid request data"
// @Failure      500 {object} map[string]string "Failed to find available room."
// @Router       /rooms/search [post]
func (h *RoomHandler) FindAvailableRoom(c *gin.Context) {
	var searchRoomRequest dto.SearchRoomRequest
	if err := c.ShouldBindJSON(&searchRoomRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}
	if !searchRoomRequest.EndDate.After(searchRoomRequest.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.start_date_must_be_before_end_date")})
		return
	}
	if searchRoomRequest.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.start_date_must_be_today_or_future")})
		return
	}
	if searchRoomRequest.MinPrice != nil && searchRoomRequest.MaxPrice != nil && *searchRoomRequest.MinPrice > *searchRoomRequest.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.min_price_must_be_less_than_max_price")})
		return
	}
	rooms, err := h.roomUseCase.SearchRoom(c.Request.Context(), &searchRoomRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.failed_to_find_available_room")})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": utils.T(c, "success.find_available_room_successful"),
		"rooms":   rooms,
	})
}
