package get_my_booking

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_service"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking"
	"log/slog"
	"net/http"
	"time"
)

// GetMyBookingsHandler godoc
// @Summary Получить список бронирований у пользователя
// @Description Получить список всех объектов бронирований с пагинацией
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id query string false "Сортировка по id. asc, desc"
// @Param booking_entity_id query string false "Сортировка по booking_entity_id. asc, desc"
// @Param start_time query string false "Сортировка по start_time. asc, desc"
// @Param end_time query string false "Сортировка по end_time. asc, desc"
// @Param status query string false "Сортировка по status. asc, desc"
// @Param user_id query string false "Сортировка по user_id. asc, desc"
// @Param search query string false "Поиск"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} bookingModels.BookingsList
// @Router /bookings/my [get]
func GetMyBookingsHandler(logger *slog.Logger, bookingDbRepo booking_db.BookingRepository, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("op", "internal/lib/services/booking_service/get_booking_by_user_id"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// берём айди пользователя, который бронирует из токена JWT
		claims := helpers.ExtractTokenClaims(ctx, log, w, r)

		// Извлекаем userId из claims
		userId := claims.AccountId

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

		response, err := booking_service.GetMyBooking(bookingDbRepo, userId, parsedQuery, logger, ctx)
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
