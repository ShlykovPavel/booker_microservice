package get_booking_type_by_id

type GetBookingTypeResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
