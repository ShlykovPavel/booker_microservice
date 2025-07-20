package get_my_booking

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

func GetMyBookingsHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_user_id"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// берём айди пользователя, который бронирует из токена JWT
		claims, ok := r.Context().Value("tokenClaims").(jwt.MapClaims)
		if !ok {
			log.Error("GetMyBookingsHandler: token claims not found in context")
			resp.RenderResponse(w, r, http.StatusUnauthorized, resp.Error("Token claims not found in context"))
			return
		}

		// Извлекаем userId из claims
		userId, ok := claims["sub"].(float64) // JWT часто возвращает числа как float64
		if !ok {
			log.Error("CreateBookingHandler: userId not found in claims")
			resp.RenderResponse(w, r, http.StatusUnauthorized, resp.Error("UserId not found in auth token"))
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

		response, err := booking_service.GetMyBooking(bookingDbRepo, int64(userId), parsedQuery, logger, ctx)
		if err != nil {
			log.Error("get my bookings failed", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		log.Debug("Success get my bookings", "response", response)
		resp.RenderResponse(w, r, http.StatusOK, response)
		return

	}
}
