package get_booking_by_time

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking"
	"github.com/ShlykovPavel/booker_microservice/models/booking/get_booking_by_time"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking/get_booking_by_time"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"time"
)

// GetBookingByTimeHandler godoc
// @Summary Получить объект бронирования по определённому промежутку времени
// @Description Получить детальную информацию о объекте бронирования
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body get_booking_by_time.GetBookingByTimeRequest true "Данные бронирования"
// @Success 200 {object} bookingModels.BookingInfo
// @Router /bookings/time [post]
func GetBookingByTimeHandler(logger *slog.Logger, bookingRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/server/booking/get_booking_by_time/get_booking_by_time_handler.go/GetBookingByTimeHandler"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		requestQuery := r.URL.Query()
		queryParser := &query_params.DefaultSortParser{
			ValidSortFields: []string{"id", "booking_entity_id", "start_time", "end_time", "status", "user_id"},
		}
		parsedQuery, err := query_params.ParseStandardQueryParams(requestQuery, log, queryParser)
		if err != nil {
			log.Error("Ошибка парсинга параметров", "error", err, "request", requestQuery)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Ошибка параметров запроса"))
			return
		}

		var getBookingDto get_booking_by_time.GetBookingByTimeRequest
		err = body.DecodeAndValidateJson(r, &getBookingDto)
		if err != nil {
			logger.Error("GetBookingByTimeHandler: error decoding body or validating", "error", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("GetBookingByTimeHandler: error decoding body", "error", err)
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

		response, err := booking_service.GetBookingByTime(bookingRepo, getBookingDto, parsedQuery, ctx, logger)
		if err != nil {
			logger.Error("GetBookingByTimeHandler", "error", err, "request", requestQuery)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		log.Info("GetBookingByTimeHandler: success", "response", response)
		resp.RenderResponse(w, r, http.StatusOK, response)
		return
	}
}
