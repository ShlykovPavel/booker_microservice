package create_booking_type

type CreateBookingTypeRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
