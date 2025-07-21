package bookingModels

import "time"

type BookingInfo struct {
	Id            int64     `json:"id"`
	UserId        int64     `json:"user_id"`
	BookingEntity int64     `json:"booking_entity"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
}

type BookingsListMetaData struct {
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}
type BookingsList struct {
	Bookings []BookingInfo        `json:"data"`
	Meta     BookingsListMetaData `json:"meta"`
}
