package create_bookingType

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/create_booking_type"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_type_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"time"
)

func CreateBookingTypeHandler(log *slog.Logger, bookingTypeRepository booking_type.BookingTypeRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("op", "internal/server/booking_type_handlers/create/create_booking_type_handler.go/CreateBookingTypeHandler"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var bookingTypeDto create_booking_type.CreateBookingTypeRequest
		err := body.DecodeAndValidateJson(r, &bookingTypeDto)
		if err != nil {
			logger.Error("CreateBookingTypeHandler: error decoding body or validating", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("CreateBookingTypeHandler: error decoding body", err)
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
		responseDto, err := booking_type_service.CreateBookingType(bookingTypeDto, bookingTypeRepository, ctx, log)
		if err != nil {
			logger.Error("CreateBookingTypeHandler: error creating booking type", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		logger.Debug("CreateBookingTypeHandler: created booking type", "response", responseDto)
		resp.RenderResponse(w, r, http.StatusCreated, responseDto)
	}
}
