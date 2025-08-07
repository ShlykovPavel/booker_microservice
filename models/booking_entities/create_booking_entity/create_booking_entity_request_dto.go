package create_booking_entity

type BookingEntity struct {
	BookingTypeID int64  `json:"booking_type_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ParentID      int64  `json:"parent_id,omitempty"`
}
