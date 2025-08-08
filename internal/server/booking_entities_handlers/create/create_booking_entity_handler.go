package create_bookingEntity_handler

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/body"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/helpers"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/booking_entities_service"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/services_models"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_entity_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type_db"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/company_db"
	"github.com/ShlykovPavel/booker_microservice/models/booking_entities/create_booking_entity"
	_ "github.com/ShlykovPavel/booker_microservice/models/booking_type/create_booking_type"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"time"
)

// CreateBookingEntityHandler godoc
// @Summary Создать объект бронирования
// @Description Создать объект бронирования
// @Tags bookingsEntity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body create_booking_entity.BookingEntity true "Данные объекта бронирования"
// @Success 201 {object} create_booking_type.ResponseId
// @Router /bookingsEntity [post]
func CreateBookingEntityHandler(log *slog.Logger, bookingTypeRepository booking_type_db.BookingTypeRepository, bookingEntityRepository booking_entity_db.BookingEntityRepository, timeout time.Duration, companyDbRepo company_db.CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("op", "internal/server/booking_entities_handlers/create/create_booking_entity_handler.go/CreateBookingEntityHandler"))

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var bookingEntityDto create_booking_entity.BookingEntity
		err := body.DecodeAndValidateJson(r, &bookingEntityDto)
		if err != nil {
			logger.Error("CreateBookingEntityHandler: error decoding body or validating", "error", err)
			if errors.Is(err, body.ErrDecodeJSON) {
				logger.Error("CreateBookingEntityHandler: error decoding body", "error", err)
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error(err.Error()))
			}
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				logger.Error("Error validating request body", "err", validationErr)
				resp.RenderResponse(w, r, http.StatusBadRequest, resp.ValidationError(validationErr))
				return
			}
			logger.Error("Unexpected error", "err", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error("internal server error"))
			return
		}
		claims := helpers.ExtractTokenClaims(ctx, log, w, r)

		createDto := services_models.CreateBookingEntityDto{
			BookingEntityInfo: bookingEntityDto,
			CompanyId:         claims.CompanyId,
			CompanyName:       claims.CompanyName,
		}

		responseDto, err := booking_entities_service.CreateBookingEntity(createDto, bookingTypeRepository, bookingEntityRepository, ctx, log, companyDbRepo)
		if err != nil {
			logger.Error("CreateBookingEntityHandler: error creating booking entity", "error", err)
			resp.RenderResponse(w, r, http.StatusInternalServerError, resp.Error(err.Error()))
			return
		}
		logger.Debug("CreateBookingEntityHandler: created booking entity", "response", responseDto)
		resp.RenderResponse(w, r, http.StatusCreated, responseDto)
	}
}
