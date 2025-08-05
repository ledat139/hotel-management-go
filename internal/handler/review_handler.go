package handler

import (
	"errors"
	"hotel-management/internal/dto"
	reviewError "hotel-management/internal/error"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewUseCase *usecase.ReviewUseCase
}

func NewReviewHandler(reviewUseCase *usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{
		reviewUseCase: reviewUseCase,
	}
}

// CreateReview godoc
// @Summary Create a review for a completed booking
// @Description Customers can create a review after checking out from a room
// @Tags Review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body dto.CreateReviewRequest true "Review content"
// @Success 201 {object} map[string]string "Review created successfully"
// @Failure 400 {object} map[string]string "Invalid request / already reviewed / not checked out"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Failed to create review"
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var createReviewRequest dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&createReviewRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, "error.invalid_request")})
		return
	}
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": utils.T(c, "error.unauthorized")})
		return
	}
	userID, ok := userIDStr.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": utils.T(c, "error.invalid_user_id")})
		return
	}
	err := h.reviewUseCase.CreateReview(c.Request.Context(), &createReviewRequest, userID)
	if err != nil {
		switch {
		case errors.Is(err, reviewError.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": utils.T(c, err.Error())})
		case errors.Is(err, reviewError.ErrBookingNotCheckedOut),
			errors.Is(err, reviewError.ErrReviewAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, err.Error())})
		case errors.Is(err, reviewError.ErrFailedToCreateReview),
			errors.Is(err, reviewError.ErrReviewCheckFailed),
			errors.Is(err, reviewError.ErrFailedToGetBooking):
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": utils.T(c, "success.review_created")})
}
