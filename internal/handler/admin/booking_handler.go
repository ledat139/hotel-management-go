package admin

import (
	"errors"
	"hotel-management/internal/constant"
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminBookingHandler struct {
	bookingUseCase *admin_usecase.BookingUseCase
}

func (h *AdminBookingHandler) ListBookings(c *gin.Context) {
	userName := c.Query("user_name")
	bookingStatus := c.Query("booking_status")

	bookings, err := h.bookingUseCase.SearchBookings(c.Request.Context(), userName, bookingStatus)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, err.Error()),
		})
		return
	}

	c.HTML(http.StatusOK, "booking.html", gin.H{
		"Title":    "admin.booking_management",
		"Bookings": bookings,
		"filters": gin.H{
			"UserName":      userName,
			"BookingStatus": bookingStatus,
		},
		"T": utils.TmplTranslateFromContext(c),
	})
}

func NewAdminBookingHandler(bookingUseCase *admin_usecase.BookingUseCase) *AdminBookingHandler {
	return &AdminBookingHandler{bookingUseCase: bookingUseCase}
}

func (h *AdminBookingHandler) GetBookingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_booking_id"),
		})
		return
	}

	booking, err := h.bookingUseCase.GetBookingDetail(c.Request.Context(), uint(id))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": utils.T(c, "error.booking_not_found"),
		})
		return
	}
	c.HTML(http.StatusOK, "booking_detail.html", gin.H{
		"Title":   "title.booking_detail",
		"Booking": booking,
		"T":       utils.TmplTranslateFromContext(c),
	})
}

func (h *AdminBookingHandler) EditBookingPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_booking_id"),
		})
		return
	}
	booking, err := h.bookingUseCase.GetBookingDetail(c.Request.Context(), uint(id))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": utils.T(c, "error.booking_not_found"),
		})
		return
	}
	c.HTML(http.StatusOK, "edit_booking.html", gin.H{
		"Title":           "title.edit_booking",
		"Booking":         booking,
		"BookingStatuses": []string{constant.BOOKED, constant.CHECKED_IN, constant.CHECKED_OUT, constant.CANCELLED, constant.NO_SHOW},
		"T":               utils.TmplTranslateFromContext(c),
	})
}

func (h *AdminBookingHandler) EditBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_booking_id"),
		})
		return
	}
	status := c.PostForm("status")
	validStatuses := map[string]bool{
		constant.BOOKED:      true,
		constant.CHECKED_IN:  true,
		constant.CHECKED_OUT: true,
		constant.CANCELLED:   true,
		constant.NO_SHOW:     true,
	}
	if status == "" || !validStatuses[status] {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_status"),
		})
		return
	}
	err = h.bookingUseCase.UpdateBookingStatus(c.Request.Context(), uint(id), status)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": utils.T(c, err.Error()),
		})
		return
	}
	c.Redirect(http.StatusSeeOther, constant.BookingManagementPath)
}
