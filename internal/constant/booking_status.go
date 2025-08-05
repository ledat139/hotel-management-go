package constant

const (
	PENDING    = "pending"
	BOOKED     = "booked"
	CHECKED_IN = "checked_in"
	CHECKED_OUT = "checked_out"
	CANCELLED = "cancelled"
	NO_SHOW   = "no_show"
)

var validBookingStatuses = map[string]bool{
	PENDING:     true,
	BOOKED:      true,
	CHECKED_IN:  true,
	CHECKED_OUT: true,
	CANCELLED:   true,
	NO_SHOW:     true,
}

func IsValidBookingStatus(status string) bool {
	return validBookingStatuses[status]
}
