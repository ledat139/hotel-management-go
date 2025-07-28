package error

import "errors"

var (
	ErrBookingNotFound      = errors.New("error.booking_not_found")
	ErrFailedToGetBooking   = errors.New("error.failed_to_get_booking")
	ErrBookingNotCheckedOut = errors.New("error.booking_not_checked_out")
	ErrReviewAlreadyExists  = errors.New("error.review_already_exists")
	ErrFailedToCreateReview = errors.New("error.failed_to_create_review")
	ErrReviewCheckFailed    = errors.New("error.review_check_failed")
)
