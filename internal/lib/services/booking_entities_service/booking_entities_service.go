package booking_entities_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_entities/create_booking_entity"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_entities/get_booking_entities_list"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_entities/get_booking_entity"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/create_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type_db"
	"log/slog"
	"strconv"
)

func CreateBookingEntity(dto create_booking_entity.BookingEntity, bookingTypeDBRepo booking_type_db.BookingTypeRepository, bookingEntityDBRepo booking_entity_db.BookingEntityRepository, ctx context.Context, log *slog.Logger) (create_booking_type.ResponseId, error) {
	log = log.With(slog.String("op", "internal/lib/services/booking_entities_service/booking_entities_service.go/CreateBookingEntity"))
	//Проверяем то тип бронирования существует
	_, err := bookingTypeDBRepo.GetBookingType(ctx, dto.BookingTypeID)
	if err != nil {
		if errors.Is(err, booking_type_db.ErrBookingTypeNotFound) {
			return create_booking_type.ResponseId{}, err
		}
		log.Error("Unexpected error while retrieve booking type", "error", err.Error())
		return create_booking_type.ResponseId{}, fmt.Errorf("failed to retrieve booking type: %w", err)

	}
	id, err := bookingEntityDBRepo.CreateBookingEntity(ctx, dto.BookingTypeID, dto.Name, dto.Description, dto.Status, dto.ParentID)
	if err != nil {
		log.Error("Ошибка создания объекта бронирования", "err", err)
		return create_booking_type.ResponseId{}, err
	}
	return create_booking_type.ResponseId{
		ID: id,
	}, nil
}

func GetBookingEntityById(id int64, bookingEntityDBRepo booking_entity_db.BookingEntityRepository, ctx context.Context, log *slog.Logger) (get_booking_entity.BookingEntityResponse, error) {
	const op = "internal/lib/services/booking_entities_service/booking_entities_service.go/GetBookingEntityById"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	BookingType, err := bookingEntityDBRepo.GetBookingEntity(ctx, id)
	if err != nil {
		if errors.Is(err, booking_entity_db.ErrBookingEntityNotFound) {
			log.Debug("Объект бронирования не найден", "err", err)
			return get_booking_entity.BookingEntityResponse{}, err
		}
		log.Error("Ошибка поиска объекта бронирования в БД", "err", err)
		return get_booking_entity.BookingEntityResponse{}, err
	}
	return get_booking_entity.BookingEntityResponse{
		Id:            BookingType.ID,
		BookingTypeID: BookingType.BookingTypeID,
		Name:          BookingType.Name,
		Description:   BookingType.Description,
		Status:        BookingType.Status,
		ParentID:      BookingType.ParentID,
	}, nil
}

func GetBookingEntitiesList(log *slog.Logger, bookingEntityDBRepo booking_entity_db.BookingEntityRepository, ctx context.Context, queryParams query_params.ListUsersParams) (get_booking_entities_list.BookingEntityList, error) {
	const op = "internal/lib/services/booking_entities_service/booking_entities_service.go/GetBookingEntitiesList"
	log = log.With(slog.String("op", op))

	result, err := bookingEntityDBRepo.GetBookingEntitiesList(ctx, queryParams.Search, queryParams.Limit, queryParams.Offset, queryParams.Sort)
	if err != nil {
		log.Error("Failed to get booking entities list", "err", err)
		return get_booking_entities_list.BookingEntityList{}, err
	}
	BookingEntitiesList := make([]get_booking_entities_list.BookingEntityInfoList, 0, len(result.BookingEntities))
	for _, bookingEntity := range result.BookingEntities {
		bookingEntityInfo := get_booking_entities_list.BookingEntityInfoList{
			Id:            bookingEntity.ID,
			BookingTypeID: bookingEntity.BookingTypeID,
			Name:          bookingEntity.Name,
			Description:   bookingEntity.Description,
			Status:        bookingEntity.Status,
			ParentID:      bookingEntity.ParentID,
		}
		BookingEntitiesList = append(BookingEntitiesList, bookingEntityInfo)
	}
	metaData := get_booking_entities_list.BookingEntityListMetaData{
		Page:   queryParams.Page,
		Total:  result.Total,
		Limit:  queryParams.Limit,
		Offset: queryParams.Offset,
	}
	bookingTEntitiesDto := get_booking_entities_list.BookingEntityList{
		BookingEntities: BookingEntitiesList,
		Meta:            metaData,
	}
	return bookingTEntitiesDto, nil
}

func UpdateBookingEntity(log *slog.Logger, bookingTypeDBRepo booking_type_db.BookingTypeRepository, bookingEntityDBRepo booking_entity_db.BookingEntityRepository, ctx context.Context, dto create_booking_entity.BookingEntity, id int64) error {
	const op = "internal/lib/services/booking_type_service/booking_type_service.go/UpdateBookingType"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	//Проверяем то тип бронирования существует
	_, err := bookingTypeDBRepo.GetBookingType(ctx, dto.BookingTypeID)
	if err != nil {
		if errors.Is(err, booking_type_db.ErrBookingTypeNotFound) {
			return err
		}
		log.Error("Unexpected error while retrieve booking entity", "error", err.Error())
		return fmt.Errorf("failed to retrieve booking entity: %w", err)
	}

	err = bookingEntityDBRepo.UpdateBookingEntity(ctx, id, dto.BookingTypeID, dto.Name, dto.Description, dto.Status, dto.ParentID)
	if err != nil {
		log.Error("Failed to update booking entity", "err", err)
		return err
	}
	return nil
}

func DeleteBookingEntity(log *slog.Logger, bookingEntityDBRepo booking_entity_db.BookingEntityRepository, ctx context.Context, id int64) error {
	const op = "internal/lib/services/booking_entities_service/booking_entities_service.go/DeleteBookingEntity"
	log = log.With(slog.String("op", op),
		slog.String("UserId", strconv.FormatInt(id, 10)))

	err := bookingEntityDBRepo.DeleteBookingEntity(ctx, id)
	if err != nil {
		log.Error("Failed to delete booking entity", "err", err)
		return err
	}
	return nil
}
