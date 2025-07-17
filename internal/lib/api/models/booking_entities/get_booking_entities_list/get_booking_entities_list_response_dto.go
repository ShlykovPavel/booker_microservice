package get_booking_entities_list

type BookingEntityInfoList struct {
	Id            int64  `json:"id"`
	BookingTypeID int64  `json:"booking_type_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ParentID      int64  `json:"parent_id,omitempty"`
}

type BookingEntityListMetaData struct {
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}
type BookingEntityList struct {
	BookingEntities []BookingEntityInfoList   `json:"data"`
	Meta            BookingEntityListMetaData `json:"meta"`
}
