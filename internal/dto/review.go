package dto

type CreateReviewRequest struct {
	BookingID uint   `json:"booking_id" binding:"required"`
	RoomID    uint   `json:"room_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}
