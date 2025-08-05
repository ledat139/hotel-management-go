package handler

import (
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUseCase *usecase.BookingUseCase
}

func NewBookingHandler(bookingUseCase *usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUseCase: bookingUseCase}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Allow a customer to create a booking with selected rooms and dates
// @Tags Booking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body dto.CreateBookingRequest true "Booking request payload"
// @Success 201 {object} dto.Response "Response with booking details and payment URL"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data/Start date must be before end date/Room is not available/Invalid client IP address"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized access"
// @Failure 404 {object} dto.ErrorResponse "Room not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to get room price/Failed to create booking/Failed to commit transaction"
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var createBookingRequest dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&createBookingRequest); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message: utils.T(c, "error.invalid_request"),
		})
		return
	}
	if !createBookingRequest.EndDate.After(createBookingRequest.StartDate) {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message: utils.T(c, "error.start_date_must_be_before_end_date"),
		})
		return
	}
	if createBookingRequest.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message: utils.T(c, "error.start_date_must_be_today_or_future"),
		})
		return
	}
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Message: utils.T(c, "error.unauthorized"),
		})
		return
	}
	userID, ok := userIDStr.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Message: utils.T(c, "error.unauthorized"),
		})
		return
	}
	clientIP := c.ClientIP()
	if clientIP == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message: utils.T(c, "error.invalid_client_ip"),
		})
		return
	}
	booking, paymentUrl, err := h.bookingUseCase.CreateBooking(c.Request.Context(), &createBookingRequest, userID, clientIP)
	if err != nil {
		switch err.Error() {
		case "error.room_not_found":
			c.JSON(http.StatusNotFound, dto.Response{
				Message: utils.T(c, "error.room_not_found"),
			})
		case "error.room_is_not_available":
			c.JSON(http.StatusBadRequest, dto.Response{
				Message: utils.T(c, "error.room_is_not_available"),
			})
		case "error.failed_to_get_room_price", "error.failed_to_create_booking", "error.failed_to_commit_transaction":
			c.JSON(http.StatusInternalServerError, dto.Response{
				Message: utils.T(c, err.Error()),
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.Response{
				Message: utils.T(c, err.Error()),
			})
		}
		return
	}
	c.JSON(http.StatusCreated, dto.Response{
		Message: utils.T(c, "success.booking_created"),
		Data: gin.H{
			"booking":     booking,
			"payment_url": paymentUrl,
		},
	})
}

// GetBookingHistory godoc
// @Summary Get booking history for current customer
// @Description Retrieve a list of past bookings for the authenticated customer
// @Tags Booking
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BookingHistoryResponse "List of booking history"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized access"
// @Failure 500 {object} dto.ErrorResponse "Failed to get booking history"
// @Router /bookings/history [get]
func (h *BookingHandler) GetBookingHistory(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Message: utils.T(c, "error.unauthorized"),
		})
		return
	}
	userID, ok := userIDStr.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Message: utils.T(c, "error.unauthorized"),
		})
		return
	}

	bookings, err := h.bookingUseCase.GetBookingHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Message: utils.T(c, "error.failed_to_get_booking_history"),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking by ID if allowed
// @Tags Booking
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string "Booking cancelled"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Failed to cancel booking"
// @Router /bookings/{id}/cancel [get]
// @Security BearerAuth
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}

	userID, exists := c.MustGet("userID").(uint)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": utils.T(c, "error.unauthorized")})
		return
	}

	err = h.bookingUseCase.CancelBooking(c.Request.Context(), uint(bookingID), userID)
	if err != nil {
		switch err.Error() {
		case "error.booking_not_found":
			c.JSON(http.StatusNotFound, gin.H{"message": utils.T(c, "error.booking_not_found")})
		case "error.failed_to_cancel_booking":
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, "error.failed_to_cancel_booking")})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.booking_cancelled")})
}
