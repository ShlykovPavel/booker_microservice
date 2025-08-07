package update_booking_type

type UpdateBookingTypeRequest struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
