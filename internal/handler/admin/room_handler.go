package admin

import (
	"fmt"
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUseCase *admin_usecase.RoomUseCase
}

func NewRoomHandler(roomUseCase *admin_usecase.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUseCase: roomUseCase}
}
func (h *RoomHandler) RoomManagementPage(c *gin.Context) {
	var query dto.RoomQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_request_data"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_management",
		})
		return
	}
	if query.HasAircon != "" && query.HasAircon != "true" && query.HasAircon != "false" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_aircon_value"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_management",
		})
		return
	}
	if query.MinPrice > 0 && query.MaxPrice > 0 && query.MinPrice > query.MaxPrice {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.min_price_greater_than_max"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_management"})
		return
	}
	rooms, err := h.roomUseCase.SearchRooms(c.Request.Context(), query)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": utils.T(c, "error.failed_to_get_room"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_management"})
		return
	}

	c.HTML(http.StatusOK, "room.html", gin.H{
		"Title": "title.room_management",
		"Rooms": rooms,
		"Query": query,
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *RoomHandler) CreateRoomPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_room.html", gin.H{
		"Title": "title.create_room",
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	formResult, err := ParseRoomForm(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_room.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.create_room",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}
	createRoomRequest := &dto.CreateRoomRequest{
		Name:          formResult.Name,
		Type:          formResult.Type,
		PricePerNight: formResult.Price,
		BedNum:        formResult.BedNum,
		HasAircon:     formResult.HasAircon,
		ViewType:      formResult.ViewType,
		Description:   formResult.Description,
		IsAvailable:   formResult.IsAvailable,
		ImageFiles:    formResult.Files,
	}

	err = h.roomUseCase.CreateRoom(c, createRoomRequest)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "create_room.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.create_room",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}

func (h *RoomHandler) RoomDetailPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_detail",
		})
		return
	}
	roomDetail, err := h.roomUseCase.GetRoomDetail(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, err.Error()),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.room_detail"})
		return
	}
	c.HTML(http.StatusOK, "room_detail.html", gin.H{
		"Title":          "title.room_detail",
		"Room":           roomDetail.Room,
		"ActiveBookings": roomDetail.ActiveBookings,
		"T":              utils.TmplTranslateFromContext(c),
	})
}

func (h *RoomHandler) EditRoomPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.edit_room",
		})
		return
	}
	room, err := h.roomUseCase.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": utils.T(c, "error.room_not_found"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.edit_room",
		})
		return
	}

	c.HTML(http.StatusOK, "edit_room.html", gin.H{
		"Title": "title.edit_room",
		"Room":  room,
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{
			"error": utils.T(c, "error.invalid_room_id"),
			"Title": "title.edit_room",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	formResult, err := ParseRoomForm(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.edit_room",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}
	// Parse deleted image id list
	deletedIDs := c.PostFormArray("delete_image_ids")

	var deletedImageIDs []int
	for _, idStr := range deletedIDs {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			deletedImageIDs = append(deletedImageIDs, id)
		}
	}

	updateReq := &dto.EditRoomRequest{
		ID:            roomID,
		Name:          formResult.Name,
		Type:          formResult.Type,
		PricePerNight: formResult.Price,
		BedNum:        formResult.BedNum,
		HasAircon:     formResult.HasAircon,
		ViewType:      formResult.ViewType,
		Description:   formResult.Description,
		IsAvailable:   formResult.IsAvailable,
		ImageDeletes:  deletedImageIDs,
	}
	fmt.Println("updateReq:", updateReq)

	err = h.roomUseCase.UpdateRoom(c, updateReq)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_room.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.edit_room",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.delete_room"})
		return
	}
	err = h.roomUseCase.DeleteRoom(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": utils.T(c, err.Error()),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.delete_room"})
		return	
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}
