package update_booking_entity

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type_db"
	"github.com/ShlykovPavel/booker_microservice/models/booking_entities/create_booking_entity"
	"github.com/ShlykovPavel/booker_microservice/models/booking_type/create_booking_type"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func UpdateBookingEntityHandler(log *slog.Logger, bookingTypeRepository booking_type_db.BookingTypeRepository, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_type_handlers/update_booking_entity/update_booking_type_handler.go/UpdateBookingEntityHandler"
		logger := log.With(slog.String("op", op))

		BookingEntityID := chi.URLParam(r, "id")
		if BookingEntityID == "" {
			log.Error("Booking Entity ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking Entity ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingEntityID, 10, 64)
		if err != nil {
			log.Error("Booking Entity ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid Booking Entity ID"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var UpdateBookingEntityDto create_booking_entity.BookingEntity
		err = body.DecodeAndValidateJson(r, &UpdateBookingEntityDto)
		if err != nil {
			logger.Error("UpdateBookingEntityHandler: error decoding body or validating", "error", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("UpdateBookingEntityHandler: error decoding body", "error", err)
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error(err.Error()))
			}
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				logger.Error("Error validating request body", "err", validationErr)
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.ValidationError(validationErr))
				return
			}
			logger.Error("Unexpected error", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("internal server error"))
			return
		}

		err = booking_entities_service.UpdateBookingEntity(log, bookingTypeRepository, bookingEntityRepository, ctx, UpdateBookingEntityDto, id)
		if err != nil {
			if errors.Is(err, booking_entity_db.ErrBookingEntityNotFound) {
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error(err.Error()))
				return
			}
			log.Error("Failed to update booking entity", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Failed updating booking entity"))
			return
		}
		log.Debug("Successfully updated booking entity", "id", id)
		resp.RenderResponse(w, r, http.StatusOK, create_booking_type.ResponseId{ID: id})

	}
}
