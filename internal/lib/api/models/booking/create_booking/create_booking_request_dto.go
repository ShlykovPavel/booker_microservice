package create_booking_dto

import "time"

type BookingRequest struct {
	UserId          int64     `json:"user_id"`
	BookingEntityId int64     `json:"booking_entity_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Status          string    `json:"status"`
	CompanyId       int64
}
