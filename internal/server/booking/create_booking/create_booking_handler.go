package create_booking

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	create_booking_dto "github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking/create_booking"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"time"
)

func CreateBookingHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration, bookingEntityDBRepo booking_entity_db.BookingEntityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/server/booking/create_booking/create_booking_handler.go/CreateBookingHandler"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		claims := helpers.ExtractTokenClaims(ctx, log, w, r)

		var createBookingDto = create_booking_dto.BookingRequest{
			UserId:    claims.AccountId,
			CompanyId: claims.CompanyId,
		}

		err := body.DecodeAndValidateJson(r, &createBookingDto)
		if err != nil {
			logger.Error("CreateBookingHandler: error decoding body or validating", "error", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("CreateBookingHandler: error decoding body", "error", err)
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
		response, err := booking_service.CreateBooking(createBookingDto, bookingDbRepo, ctx, logger, bookingEntityDBRepo)
		if err != nil {
			logger.Error("CreateBookingHandler: error creating booking", "error", err)
			if errors.Is(err, booking_service.ErrBookingNotAvailable) {
				logger.Warn("CreateBookingHandler: booking not available", "error", err)
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking not available"))
				return
			}
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		resp.RenderResponse(w, r, http.StatusOK, response)
		return
	}

}
