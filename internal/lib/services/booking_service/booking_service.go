package booking_service

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/models/booking"
	"github.com/ShlykovPavel/booker_microservice/models/booking/create_booking"
	"github.com/ShlykovPavel/booker_microservice/models/booking/get_booking_by_time"
	"github.com/ShlykovPavel/booker_microservice/models/booking_type/create_booking_type"
	"log/slog"
)

var ErrBookingNotAvailable = errors.New("Booking not available")

func CreateBooking(dto create_booking_dto.BookingRequest, bookingRepo booking_db.BookingRepository, ctx context.Context, log *slog.Logger, bookingEntityDBRepo booking_entity_db.BookingEntityRepository) (create_booking_type.ResponseId, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/booking_service.go/CreateBooking"))

	//Проверка, что время свободно
	available, err := bookingRepo.CheckBookingAvailability(ctx, dto.BookingEntityId, dto.StartTime, dto.EndTime)
	if err != nil {
		log.Error("Check Booking Availability failed", "error", err.Error())
		return create_booking_type.ResponseId{}, err
	}
	if !available {
		return create_booking_type.ResponseId{}, ErrBookingNotAvailable
	}

	entityInfo, err := bookingEntityDBRepo.GetBookingEntity(ctx, dto.BookingEntityId)
	if err != nil {
		log.Error("Get booking entity failed", "error", err.Error())
		return create_booking_type.ResponseId{}, err
	}

	// Проверяем, что companyId из токена совпадает с companyId сущности
	if entityInfo.CompanyId != dto.CompanyId {
		log.Error("Company ID mismatch", "tokenCompanyId", dto.CompanyId, "entityCompanyId", entityInfo.CompanyId)
		return create_booking_type.ResponseId{}, errors.New("access denied: company ID mismatch")
	}

	id, err := bookingRepo.CreateBooking(ctx, dto)
	if err != nil {
		log.Error("CreateBooking failed", "error", err)
		return create_booking_type.ResponseId{}, err
	}
	return create_booking_type.ResponseId{ID: id}, nil
}

// GetBookingByTime получить все бронирования за определённый промежуток времени
func GetBookingByTime(bookingRepo booking_db.BookingRepository, dto get_booking_by_time.GetBookingByTimeRequest, queryParams query_params.ListQueryParams, ctx context.Context, log *slog.Logger) ([]bookingModels.BookingInfo, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_time"))

	bookings, err := bookingRepo.GetBookingsByTime(ctx, dto.StartTime, dto.EndTime, queryParams)
	if err != nil {
		log.Error("GetBookingByTime failed", "error", err)
		return []bookingModels.BookingInfo{}, err
	}
	bookingsList := make([]bookingModels.BookingInfo, 0, len(bookings))
	for _, booking := range bookings {
		bookingInfo := bookingModels.BookingInfo{
			Id:            booking.Id,
			UserId:        booking.UserId,
			BookingEntity: booking.BookingEntityId,
			Status:        booking.Status,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
		}
		bookingsList = append(bookingsList, bookingInfo)
	}
	return bookingsList, nil
}

// GetMyBooking Получить все бронирования у выбранного пользователя
func GetMyBooking(bookingRepo booking_db.BookingRepository, userId int64, queryParams query_params.ListQueryParams, log *slog.Logger, ctx context.Context) (bookingModels.BookingsList, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_user_id"))

	bookings, err := bookingRepo.GetBookingsByUserId(ctx, userId, queryParams)
	if err != nil {
		log.Error("GetBookingByUserId failed", "error", err)
		return bookingModels.BookingsList{}, err
	}

	bookingsList := make([]bookingModels.BookingInfo, 0, len(bookings.Bookings))
	for _, booking := range bookings.Bookings {
		bookingInfo := bookingModels.BookingInfo{
			Id:            booking.Id,
			UserId:        booking.UserId,
			BookingEntity: booking.BookingEntityId,
			Status:        booking.Status,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
		}
		bookingsList = append(bookingsList, bookingInfo)
	}
	metaData := bookingModels.BookingsListMetaData{
		Page:   queryParams.Page,
		Limit:  queryParams.Limit,
		Total:  bookings.Total,
		Offset: queryParams.Offset,
	}
	bookingListDto := bookingModels.BookingsList{
		Bookings: bookingsList,
		Meta:     metaData,
	}
	return bookingListDto, nil
}

// GetBookingByBookingEntity Получить все бронирования у выбранной сущности бронирования
func GetBookingByBookingEntity(bookingRepo booking_db.BookingRepository, bookingEntityId int64, queryParams query_params.ListQueryParams, log *slog.Logger, ctx context.Context) (bookingModels.BookingsList, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_user_id"))

	bookings, err := bookingRepo.GetBookingsByBookingEntity(ctx, bookingEntityId, queryParams)
	if err != nil {
		log.Error("GetBookingsByBookingEntity failed", "error", err)
		return bookingModels.BookingsList{}, err
	}

	bookingsList := make([]bookingModels.BookingInfo, 0, len(bookings.Bookings))
	for _, booking := range bookings.Bookings {
		bookingInfo := bookingModels.BookingInfo{
			Id:            booking.Id,
			UserId:        booking.UserId,
			BookingEntity: booking.BookingEntityId,
			Status:        booking.Status,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
		}
		bookingsList = append(bookingsList, bookingInfo)
	}
	metaData := bookingModels.BookingsListMetaData{
		Page:   queryParams.Page,
		Limit:  queryParams.Limit,
		Total:  bookings.Total,
		Offset: queryParams.Offset,
	}
	bookingListDto := bookingModels.BookingsList{
		Bookings: bookingsList,
		Meta:     metaData,
	}
	return bookingListDto, nil
}

// GetBookingById получение информации по определённому бронированию
func GetBookingById(bookingRepo booking_db.BookingRepository, bookingId int64, log *slog.Logger, ctx context.Context) (bookingModels.BookingInfo, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_id"))

	booking, err := bookingRepo.GetBookingById(ctx, bookingId)
	if err != nil {
		log.Error("GetBookingById failed", "error", err)
		return bookingModels.BookingInfo{}, err
	}
	return bookingModels.BookingInfo{
		Id:            booking.Id,
		UserId:        booking.UserId,
		BookingEntity: booking.BookingEntityId,
		Status:        booking.Status,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
	}, nil
}

// UpdateBooking обновление бронирования
func UpdateBooking(bookingRepo booking_db.BookingRepository, dto create_booking_dto.BookingRequest, bookingId int64, log *slog.Logger, ctx context.Context) (create_booking_type.ResponseId, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_service/update_booking"))

	//Проверка, что время свободно
	available, err := bookingRepo.CheckBookingAvailability(ctx, dto.BookingEntityId, dto.StartTime, dto.EndTime, bookingId)
	if err != nil {
		log.Error("Check Booking Availability failed", "error", err.Error())
		return create_booking_type.ResponseId{}, err
	}
	if !available {
		return create_booking_type.ResponseId{}, ErrBookingNotAvailable
	}

	updateDbDto := booking_db.BookingInfo{
		Id:              bookingId,
		UserId:          dto.UserId,
		BookingEntityId: dto.BookingEntityId,
		Status:          dto.Status,
		StartTime:       dto.StartTime,
		EndTime:         dto.EndTime,
	}

	err = bookingRepo.UpdateBooking(ctx, updateDbDto, bookingId)
	if err != nil {
		log.Error("Update Booking failed", "error", err)
		return create_booking_type.ResponseId{}, err
	}

	return create_booking_type.ResponseId{ID: bookingId}, nil
}

func DeleteBooking(bookingRepo booking_db.BookingRepository, bookingId int64, log *slog.Logger, ctx context.Context) error {

	log = log.With(slog.String("op", "internal/lib/services/booking_service/delete_booking"))

	err := bookingRepo.DeleteBooking(ctx, bookingId)
	if err != nil {
		log.Error("DeleteBooking failed", "error", err)
		return err
	}

	return nil
}
