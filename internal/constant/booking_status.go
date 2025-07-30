package constant

const (
	CHECKED_OUT = "checked_out"
	BOOKED      = "booked"
	CANCELLED   = "cancelled"
	CHECKED_IN  = "checked_in"
	NO_SHOW     = "no_show"
)

var validBookingStatuses = map[string]bool{
	BOOKED:      true,
	CHECKED_IN:  true,
	CHECKED_OUT: true,
	CANCELLED:   true,
	NO_SHOW:     true,
}

func IsValidBookingStatus(status string) bool {
	return validBookingStatuses[status]
}
