package update_booking

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/ShlykovPavel/booker_microservice/models/booking/create_booking"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking/create_booking"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking_type/create_booking_type"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// UpdateBookingHandler godoc
// @Summary Обновить бронирование
// @Description Обновить бронирование
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID бронирования"
// @Param input body create_booking_dto.BookingRequest true "Данные бронирования"
// @Success 200 {object} create_booking_type.ResponseId
// @Router /bookings/{id} [put]
func UpdateBookingHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/update_booking"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		claims := helpers.ExtractTokenClaims(ctx, log, w, r)

		BookingID := chi.URLParam(r, "id")
		if BookingID == "" {
			log.Error("Booking  ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingID, 10, 64)
		if err != nil {
			log.Error("Booking ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid user ID"))
			return
		}

		var updateBookingDto = create_booking_dto.BookingRequest{
			UserId:    claims.AccountId,
			CompanyId: claims.CompanyId,
		}
		err = body.DecodeAndValidateJson(r, &updateBookingDto)
		if err != nil {
			logger.Error("UpdateBookingHandler: error decoding body or validating", "error", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("UpdateBookingHandler: error decoding body", "error", err)
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
		logger.Debug("UpdateBookingHandler: parsed body from request", "body", updateBookingDto)
		response, err := booking_service.UpdateBooking(bookingDbRepo, updateBookingDto, id, logger, ctx)
		if err != nil {
			logger.Error("UpdateBookingHandler", "error", err)
			if errors.Is(err, booking_service.ErrBookingNotAvailable) {
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking not available"))
			}
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		logger.Info("Success UpdateBookingHandler", "response", response)
		resp.RenderResponse(w, r, http.StatusOK, response)
		return

	}
}
