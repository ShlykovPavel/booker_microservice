package get_booking_entity

type BookingEntityResponse struct {
	Id            int64  `json:"id"`
	BookingTypeID int64  `json:"booking_type_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ParentID      int64  `json:"parent_id,omitempty"`
}
