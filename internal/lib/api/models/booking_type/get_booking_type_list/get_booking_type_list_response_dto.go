package get_booking_type_list

type BookingTypeInfoList struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BookingTypeListMetaData struct {
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}
type BookingTypeList struct {
	BookingTypes []BookingTypeInfoList   `json:"data"`
	Meta         BookingTypeListMetaData `json:"meta"`
}
