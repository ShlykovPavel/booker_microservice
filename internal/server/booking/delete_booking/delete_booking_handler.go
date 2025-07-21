package delete_booking

import (
	"context"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func DeleteBookingHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/delete_booking"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

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

		err = booking_service.DeleteBooking(bookingDbRepo, id, logger, ctx)
		if err != nil {
			log.Error("DeleteBooking failed", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		resp.RenderResponse(w, r, http.StatusNoContent, nil)
		return
	}
}
