package booking_type_service

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/create_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	"log/slog"
)

func CreateBookingType(dto create_booking_type.CreateBookingTypeRequest, bookingTypeDBRepo booking_type.BookingTypeRepository, ctx context.Context, log *slog.Logger) (create_booking_type.CreateBookingTypeResponse, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_type_service/booking_type_service.go/CreateBookingType"))

	id, err := bookingTypeDBRepo.CreateBookingType(ctx, dto.Name, dto.Description)
	if err != nil {
		log.Error("Ошибка создания типа бронирования", "err", err)
		return create_booking_type.CreateBookingTypeResponse{}, err
	}
	return create_booking_type.CreateBookingTypeResponse{
		ID: id,
	}, nil
}
