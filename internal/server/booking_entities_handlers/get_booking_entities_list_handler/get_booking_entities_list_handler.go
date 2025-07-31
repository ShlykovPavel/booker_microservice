package get_booking_entities_list_handler

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/services_models"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/company_db"
	"log/slog"
	"net/http"
	"time"
)

func GetBookingEntitiesListHandler(logger *slog.Logger, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration, companyDbRepo company_db.CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/server/booking_entities_handlers/get_booking_entities_list_handler/get_booking_entities_list_handler.go/GetBookingEntitiesListHandler"
		log := logger.With(slog.String("op", op))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		requestQuery := r.URL.Query()

		queryParser := &query_params.DefaultSortParser{
			ValidSortFields: []string{"id", "booking_type_id", "name", "description", "status"},
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

		entitiesList, err := booking_entities_service.GetBookingEntitiesList(log, bookingEntityRepository, ctx, parsedQuery, companyDto, companyDbRepo)
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
