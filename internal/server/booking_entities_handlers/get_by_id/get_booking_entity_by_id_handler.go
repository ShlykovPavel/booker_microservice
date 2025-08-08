package get_bookingEntity_by_id_handler

import (
	"context"
	"errors"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking_entities/get_booking_entity"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// GetBookingEntityByIdHandler godoc
// @Summary Получить объект бронирования по ID
// @Description Получить детальную информацию о объекте бронирования
// @Tags bookingsEntity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID объекта бронирования"
// @Success 200 {object} get_booking_entity.BookingEntityResponse
// @Router /bookingsEntity/{id} [get]
func GetBookingEntityByIdHandler(log *slog.Logger, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("op", "internal/server/booking_type_handlers/get_by_id/get_booking_entity_by_id_handler.go/GetBookingEntityByIdHandler"))

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
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		bookingEntityInfo, err := booking_entities_service.GetBookingEntityById(id, bookingEntityRepository, ctx, logger)
		if err != nil {
			if errors.Is(err, booking_entity_db.ErrBookingEntityNotFound) {
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error("Booking entity not found"))
				return
			}
			log.Error("Error while getting booking entity by id", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Something went wrong, while getting booking entity"))
			return
		}
		log.Debug("Successful get booking entity by id", "user", bookingEntityInfo)
		resp.RenderResponse(w, r, http.StatusOK, bookingEntityInfo)
		return
	}
}
