package services_models

import "github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_entities/create_booking_entity"

type CreateBookingEntityDto struct {
	BookingEntityInfo create_booking_entity.BookingEntity
	CompanyId         int64
	CompanyName       string
}
