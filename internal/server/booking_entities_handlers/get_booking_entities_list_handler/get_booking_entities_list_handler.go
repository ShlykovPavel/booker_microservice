package get_booking_entities_list_handler

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"log/slog"
	"net/http"
	"time"
)

func GetBookingEntitiesListHandler(logger *slog.Logger, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_entities_handlers/get_booking_entities_list_handler/get_booking_entities_list_handler.go/GetBookingEntitiesListHandler"
		log := logger.With(slog.String("op", op))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		requestQuery := r.URL.Query()
		parsedQuery, err := query_params.ParseStandardQueryParams(requestQuery, log)
		if err != nil {
			log.Error("Ошибка парсинга параметров", "error", err, "request", requestQuery)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Ошибка параметров запроса"))
			return
		}

		entitiesList, err := booking_entities_service.GetBookingEntitiesList(log, bookingEntityRepository, ctx, parsedQuery)
		if err != nil {
			log.Error("Error while getting booking entity list", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Something went wrong, while getting booking entity list"))
			return
		}
		log.Debug("Successful get booking entity list")
		resp.RenderResponse(w, r, http.StatusOK, entitiesList)
		return

	}
}
