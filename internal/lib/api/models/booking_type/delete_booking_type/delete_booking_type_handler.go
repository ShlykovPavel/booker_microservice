package delete_booking_type

import (
	"context"
	"errors"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_type_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func DeleteBookingTypeHandler(logger *slog.Logger, bookingTypeRepository booking_type.BookingTypeRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/api/models/booking_type/delete_booking_type/delete_booking_type_handler.go/DeleteBookingTypeHandler"))
		bookingTypeID := chi.URLParam(r, "id")
		if bookingTypeID == "" {
			log.Error("bookingType ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("bookingTypeID is required"))
			return
		}
		id, err := strconv.ParseInt(bookingTypeID, 10, 64)
		if err != nil {
			log.Error("bookingTypeID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid bookingTypeID"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		err = booking_type_service.DeleteBookingType(logger, bookingTypeRepository, ctx, id)
		if err != nil {
			if errors.Is(err, booking_type.ErrBookingTypeNotFound) {
				log.Error("Booking type not found", "error", err)
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error("Booking type not found"))
				return
			}
			log.Error("Error deleting Booking type", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Error deleting Booking type"))
			return
		}
		log.Info("Deleted Booking type", "Booking typeID", bookingTypeID)
		resp.RenderResponse(w, r, http.StatusNoContent, nil)
	}

}
