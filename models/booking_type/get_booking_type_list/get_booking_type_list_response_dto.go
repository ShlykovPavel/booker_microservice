package get_booking_type_list

// BookingTypeInfoList represents a single booking type
type BookingTypeInfoList struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// BookingTypeListMetaData represents response metadata
type BookingTypeListMetaData struct {
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

// BookingTypeList represents response
type BookingTypeList struct {
	BookingTypes []BookingTypeInfoList   `json:"data"`
	Meta         BookingTypeListMetaData `json:"meta"`
}
