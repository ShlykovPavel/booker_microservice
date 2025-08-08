package get_entities_by_booking_type

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/services_models"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/company_db"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking_entities/get_booking_type_entities"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// GetBookingEntitiesListHandler godoc
// @Summary Получить список объектов бронирований у определённого типа бронирования
// @Description Получить список всех объектов бронирований с пагинацией
// @Tags bookingsEntity
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID типа бронирования"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} get_booking_type_entities.BookingTypeEntitiesResponse
// @Router /bookingsType/{id}/bookingsEntity [get]
func GetBookingEntitiesListHandler(logger *slog.Logger, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration, companyDbRepo company_db.CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With("op", "get_entities_by_booking_type.GetBookingEntitiesListHandler")

		bookingTypeID := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(bookingTypeID, 10, 64)
		if err != nil {
			log.Error("Booking Type ID is invalid", "error", err)
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid Booking Type ID"))
			return
		}
		if bookingTypeID == "" {
			log.Error("Booking Type ID is empty")
			resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Booking Type ID is required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		claims := helpers.ExtractTokenClaims(ctx, log, w, r)
		companyDto := services_models.CompanyInfo{
			CompanyId:   claims.CompanyId,
			CompanyName: claims.CompanyName,
		}

		entitiesList, err := booking_entities_service.GetEntitiesByType(logger, bookingEntityRepository, ctx, id, companyDto, companyDbRepo)
		if err != nil {
			log.Error("Error while getting booking entity list", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("Something went wrong, while getting booking entity list"))
			return
		}
		log.Debug("Successful get booking entity list by booking type", "bookingType", bookingTypeID)
		resp.RenderResponse(w, r, http.StatusOK, entitiesList)
		return
	}
}
