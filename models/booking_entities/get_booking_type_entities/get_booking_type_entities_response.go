package get_booking_type_entities

type BookingTypeEntitiesResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
