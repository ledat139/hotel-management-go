package error

import "errors"

var (
	ErrBookingNotFound         = errors.New("error.booking_not_found")
	ErrFailedToGetBooking      = errors.New("error.failed_to_get_booking")
	ErrBookingNotCheckedOut    = errors.New("error.booking_not_checked_out")
	ErrReviewAlreadyExists     = errors.New("error.review_already_exists")
	ErrFailedToCreateReview    = errors.New("error.failed_to_create_review")
	ErrReviewCheckFailed       = errors.New("error.review_check_failed")
	ErrBookingHasPaid          = errors.New("error.booking_has_paid")
	ErrPaymentNotFound         = errors.New("error.payment_not_found")
	ErrFailedToGetPayment      = errors.New("error.failed_to_get_payment")
	ErrPaymentAlreadyProcessed = errors.New("error.payment_already_processed")
	ErrFailedToUpdatePayment   = errors.New("error.failed_to_update_payment")
	ErrFailedToUpdateBooking   = errors.New("error.failed_to_update_booking")
	ErrFailedToCreateBill      = errors.New("error.failed_to_create_bill")
)
