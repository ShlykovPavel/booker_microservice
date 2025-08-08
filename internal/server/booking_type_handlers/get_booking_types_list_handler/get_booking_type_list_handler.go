package get_booking_types_list_handler

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_type_service"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/services_models"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/company_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking_type/get_booking_type_list"
	"log/slog"
	"net/http"
	"time"
)

// GetBookingTypesListHandler godoc
// @Summary Получить список типов бронирований
// @Description Получить список всех типов бронирований с пагинацией
// @Tags bookingsType
// @Produce json
// @Security BearerAuth
// @Param id query string false "Сортировка по id. asc, desc"
// @Param name query string false "Сортировка по name. asc, desc"
// @Param description query string false "Сортировка по description. asc, desc"
// @Param search query string false "Поиск"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} get_booking_type_list.BookingTypeList
// @Router /bookingsType [get]
func GetBookingTypesListHandler(logger *slog.Logger, bookingTypeRepository booking_type_db.BookingTypeRepository, timeout time.Duration, companyDbRepo company_db.CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_type_handlers/get_booking_types_list_handler/get_booking_type_list_handler.go/GetBookingTypesListHandler"
		log := logger.With(slog.String("op", op))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		requestQuery := r.URL.Query()

		queryParser := &query_params.DefaultSortParser{
			ValidSortFields: []string{"id", "name", "description"},
		}
		parsedQuery, err := query_params.ParseStandardQueryParams(requestQuery, log, queryParser)
		if err != nil {
			log.Error("Ошибка парсинга параметров", "error", err, "request", requestQuery)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Ошибка параметров запроса"))
			return
		}
		claims := helpers.ExtractTokenClaims(ctx, log, w, r)
		companyDto := services_models.CompanyInfo{
			CompanyId:   claims.CompanyId,
			CompanyName: claims.CompanyName,
		}
		userList, err := booking_type_service.GetBookingTypeList(log, bookingTypeRepository, ctx, parsedQuery, companyDto, companyDbRepo)
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
