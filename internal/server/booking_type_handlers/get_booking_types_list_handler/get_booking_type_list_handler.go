package get_booking_types_list_handler

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_type_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type_db"
	"log/slog"
	"net/http"
	"time"
)

func GetBookingTypesListHandler(logger *slog.Logger, bookingTypeRepository booking_type_db.BookingTypeRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_type_handlers/get_booking_types_list_handler/get_booking_type_list_handler.go/GetBookingTypesListHandler"
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

		userList, err := booking_type_service.GetBookingTypeList(log, bookingTypeRepository, ctx, parsedQuery)
		if err != nil {
			log.Error("Error while getting booking type list", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Something went wrong, while getting booking type list"))
			return
		}
		log.Debug("Successful get booking type list")
		resp.RenderResponse(w, r, http.StatusOK, userList)
		return

	}
}
