package get_booking_by_time

import "time"

type GetBookingByTimeRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
