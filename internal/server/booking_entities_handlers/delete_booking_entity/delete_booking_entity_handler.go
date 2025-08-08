package delete_booking_entity

import (
	"context"
	"errors"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// DeleteBookingEntityHandler godoc
// @Summary Удалить объект бронирования
// @Description Удалить объект бронирования по ID
// @Tags bookingsEntity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID объекта бронирования"
// @Success 204
// @Router /bookingsEntity/{id} [delete]
func DeleteBookingEntityHandler(logger *slog.Logger, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/api/models/booking_type_db/delete_booking_entity/delete_booking_entity_handler.go/DeleteBookingEntityHandler"))
		bookingEntityID := chi.URLParam(r, "id")
		if bookingEntityID == "" {
			log.Error("BookingEntity ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("BookingEntityID is required"))
			return
		}
		id, err := strconv.ParseInt(bookingEntityID, 10, 64)
		if err != nil {
			log.Error("bookingEntityID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid bookingEntityID"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		err = booking_entities_service.DeleteBookingEntity(logger, bookingEntityRepository, ctx, id)
		if err != nil {
			if errors.Is(err, booking_entity_db.ErrBookingEntityNotFound) {
				log.Error("Booking entity not found", "error", err)
				resp.RenderResponse(w, r, http.StatusNotFound, resp.Error("Booking entity not found"))
				return
			}
			log.Error("Error deleting Booking entity", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Error deleting Booking entity"))
			return
		}
		log.Info("Deleted Booking entity", "Booking entityID", bookingEntityID)
		resp.RenderResponse(w, r, http.StatusNoContent, nil)
	}

}
