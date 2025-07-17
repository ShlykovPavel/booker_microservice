package get_bookingType_by_id_handler

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

func GetBookingTypeByIdHandler(log *slog.Logger, bookingTypeRepository booking_type.BookingTypeRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("op", "internal/server/booking_type_handlers/get_by_id/get_booking_type_by_id_handler.go/GetBookingTypeByIdHandler"))

		BookingTypeID := chi.URLParam(r, "id")
		if BookingTypeID == "" {
			log.Error("User ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("User ID is required"))
			return
		}
		id, err := strconv.ParseInt(BookingTypeID, 10, 64)
		if err != nil {
			log.Error("User ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid user ID"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		bookingTypeInfo, err := booking_type_service.GetBookingTypeById(id, bookingTypeRepository, ctx, logger)
		if err != nil {
			if errors.Is(err, booking_type.ErrBookingTypeNotFound) {
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error("Booking type not found"))
				return
			}
			log.Error("Error while getting booking type by id", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Something went wrong, while getting booking type"))
			return
		}
		log.Debug("Successful get user by id", "user", bookingTypeInfo)
		resp.RenderResponse(w, r, http.StatusOK, bookingTypeInfo)
		return
	}
}
