package update_booking_type

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/create_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/update_booking_type"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_type_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func UpdateBookingTypeHandler(log *slog.Logger, bookingTypeRepository booking_type.BookingTypeRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_type_handlers/update_booking_type/update_booking_type_handler.go/UpdateBookingTypeHandler"
		logger := log.With(slog.String("op", op))

		BookingTypeID := chi.URLParam(r, "id")
		if BookingTypeID == "" {
			log.Error("User ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("BookingType ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingTypeID, 10, 64)
		if err != nil {
			log.Error("Booking Type ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid BookingType ID"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var UpdateBookingTypeDto update_booking_type.UpdateBookingTypeRequest
		err = body.DecodeAndValidateJson(r, &UpdateBookingTypeDto)
		if err != nil {
			logger.Error("UpdateBookingTypeHandler: error decoding body or validating", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("UpdateBookingTypeHandler: error decoding body", err)
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

		err = booking_type_service.UpdateBookingType(log, bookingTypeRepository, ctx, UpdateBookingTypeDto, id)
		if err != nil {
			if errors.Is(err, booking_type.ErrBookingTypeNotFound) {
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error(err.Error()))
				return
			}
			log.Error("Failed to update booking type", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Failed updating booking type"))
			return
		}
		log.Debug("Successfully updated booking type", "id", id)
		resp.RenderResponse(w, r, http.StatusOK, create_booking_type.CreateBookingTypeResponse{ID: id})

	}
}
