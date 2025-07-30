package company_service

import (
	"context"
	"errors"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/company_db"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Company struct {
	CompanyId   int64  `json:"company_id" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
}

func CreateCompany(companyDbRepo company_db.CompanyRepository, logger *slog.Logger, ctx context.Context, companyId int64, companyName string) error {
	log := logger.With(slog.String("op", "company_service/CreateCompany"))

	_, err := companyDbRepo.CreateCompany(ctx, companyId, companyName)
	if err != nil {
		log.Error("Failed create company", "err", err)
		return err
	}
	return nil
}

func CheckCompany(companyDbRepo company_db.CompanyRepository, logger *slog.Logger, ctx context.Context, companyId int64, companyName string) (bool, error) {
	log := logger.With(slog.String("op", "company_service/CheckCompany"))

	// Проверяем наличие компании
	id, err := companyDbRepo.GetCompany(ctx, companyId)
	if err != nil {
		// Если компания не найдена, создаем новую
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Company not found, creating new one", "company_id", companyId)
			_, createErr := companyDbRepo.CreateCompany(ctx, companyId, companyName)
			if createErr != nil {
				log.Error("Failed to create company", "err", createErr)
				return false, createErr
			}
			log.Debug("Company created successfully", "company_id", companyId)
			return true, nil
		}
		// Если другая ошибка, возвращаем её
		log.Error("Failed to get company", "err", err)
		return false, err
	}

	// Если id валидный, компания существует
	if id > 0 {
		log.Debug("Company found", "company_id", companyId)
		return true, nil
	}
	return false, nil
}
