package get_booking_by_id

import (
	"context"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// GetBookingByIdHandler godoc
// @Summary Получить бронирование по ID
// @Description Получить детальную информацию о бронировании
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID бронирования"
// @Success 200 {object} bookingModels.BookingInfo
// @Router /bookings/{id} [get]
func GetBookingByIdHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_id"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		BookingID := chi.URLParam(r, "id")
		if BookingID == "" {
			log.Error("Booking  ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking Entity ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingID, 10, 64)
		if err != nil {
			log.Error("Booking Entity ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid user ID"))
			return
		}
		response, err := booking_service.GetBookingById(bookingDbRepo, id, logger, ctx)
		if err != nil {
			log.Error("failed to get booking by id", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}

		log.Info("Successfully fetched booking by id", "id", BookingID)
		resp.RenderResponse(w, r, http.StatusOK, response)
		return

	}
}
