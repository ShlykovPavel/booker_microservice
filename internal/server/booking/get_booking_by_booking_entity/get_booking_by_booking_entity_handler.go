package get_booking_by_booking_entity

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func GetMyBookingsHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_user_id"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		BookingEntityID := chi.URLParam(r, "id")
		if BookingEntityID == "" {
			log.Error("Booking Entity ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking Entity ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingEntityID, 10, 64)
		if err != nil {
			log.Error("Booking Entity ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid user ID"))
			return
		}

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

		response, err := booking_service.GetBookingByBookingEntity(bookingDbRepo, id, parsedQuery, logger, ctx)
		if err != nil {
			log.Error("get bookings for booking entity failed", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}

		log.Debug("Success get bookings for booking entity", "response", response)
		resp.RenderResponse(w, r, http.StatusOK, response)
		return

	}
}
