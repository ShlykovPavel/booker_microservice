package booking_type_service

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/create_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/get_booking_type_by_id"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/get_booking_type_list"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/update_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	"log/slog"
	"strconv"
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

func GetBookingTypeById(id int64, bookingTypeDBRepo booking_type.BookingTypeRepository, ctx context.Context, log *slog.Logger) (get_booking_type_by_id.GetBookingTypeResponse, error) {
	const op = "internal/lib/services/booking_type_service/booking_type_service.go/GetBookingTypeById"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	BookingType, err := bookingTypeDBRepo.GetBookingType(ctx, id)
	if err != nil {
		if errors.Is(err, booking_type.ErrBookingTypeNotFound) {
			log.Debug("Тип бронирования не найден", "err", err)
			return get_booking_type_by_id.GetBookingTypeResponse{}, err
		}
		log.Error("Ошибка поиска типа бронирования в БД", "err", err)
		return get_booking_type_by_id.GetBookingTypeResponse{}, err
	}
	return get_booking_type_by_id.GetBookingTypeResponse{
		Id:          id,
		Name:        BookingType.Name,
		Description: BookingType.Description,
	}, nil
}

func GetBookingTypeList(log *slog.Logger, bookingTypeDBRepo booking_type.BookingTypeRepository, ctx context.Context, queryParams query_params.ListUsersParams) (get_booking_type_list.BookingTypeList, error) {
	const op = "internal/lib/services/user_service/user_service.go/GetUserList"
	log = log.With(slog.String("op", op))

	result, err := bookingTypeDBRepo.GetBookingTypeList(ctx, queryParams.Search, queryParams.Limit, queryParams.Offset, queryParams.Sort)
	if err != nil {
		log.Error("Failed to get users list", "err", err)
		return get_booking_type_list.BookingTypeList{}, err
	}
	BookingTypeList := make([]get_booking_type_list.BookingTypeInfoList, 0, len(result.BookingTypes))
	for _, bookingType := range result.BookingTypes {
		bookingTypeInfo := get_booking_type_list.BookingTypeInfoList{
			Id:          bookingType.ID,
			Name:        bookingType.Name,
			Description: bookingType.Description,
		}
		BookingTypeList = append(BookingTypeList, bookingTypeInfo)
	}
	metaData := get_booking_type_list.BookingTypeListMetaData{
		Page:   queryParams.Page,
		Total:  result.Total,
		Limit:  queryParams.Limit,
		Offset: queryParams.Offset,
	}
	bookingTypesDto := get_booking_type_list.BookingTypeList{
		BookingTypes: BookingTypeList,
		Meta:         metaData,
	}
	return bookingTypesDto, nil

}

func UpdateBookingType(log *slog.Logger, bookingTypeDBRepo booking_type.BookingTypeRepository, ctx context.Context, dto update_booking_type.UpdateBookingTypeRequest, id int64) error {
	const op = "internal/lib/services/booking_type_service/booking_type_service.go/UpdateBookingType"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	err := bookingTypeDBRepo.UpdateBookingType(ctx, id, dto.Name, dto.Description)
	if err != nil {
		log.Error("Failed to update booking type", "err", err)
		return err
	}
	return nil
}

func DeleteBookingType(log *slog.Logger, bookingTypeDBRepo booking_type.BookingTypeRepository, ctx context.Context, id int64) error {
	const op = "internal/lib/services/booking_type_service/booking_type_service.go/DeleteBookingType"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	err := bookingTypeDBRepo.DeleteBookingType(ctx, id)
	if err != nil {
		log.Error("Failed to delete booking type", "err", err)
		return err
	}
	return nil
}
