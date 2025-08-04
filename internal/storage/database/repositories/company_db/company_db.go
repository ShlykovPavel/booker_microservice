package company_db

import (
	"context"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type CompanyRepositoryImpl struct {
	dbPoll *pgxpool.Pool
	log    *slog.Logger
}

func NewCompanyRepository(db *pgxpool.Pool, log *slog.Logger) *CompanyRepositoryImpl {
	return &CompanyRepositoryImpl{
		dbPoll: db,
		log:    log,
	}
}

type CompanyRepository interface {
	CreateCompany(ctx context.Context, companyId int64, companyName string) (int64, error)
	GetCompany(ctx context.Context, companyId int64) (int64, error)
	DeleteCompany(ctx context.Context, companyId int64) error
}

func (c *CompanyRepositoryImpl) CreateCompany(ctx context.Context, companyId int64, companyName string) (int64, error) {
	query := `INSERT INTO companies (company_id, company_name) VALUES ($1, $2) RETURNING id`
	var id int64
	err := c.dbPoll.QueryRow(ctx, query, companyId, companyName).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		c.log.Error("Failed to create company", "error", err)
		return 0, database.PsqlErrorHandler(dbErr)
	}
	c.log.Debug("Created company", "id", id)
	return id, nil
}
func (c *CompanyRepositoryImpl) GetCompany(ctx context.Context, companyId int64) (int64, error) {
	query := `SELECT id FROM companies WHERE company_id = $1`
	var id int64
	err := c.dbPoll.QueryRow(ctx, query, companyId).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		c.log.Error("Failed to get company", "error", err)
		return 0, database.PsqlErrorHandler(dbErr)
	}
	c.log.Debug("Get company in db", "id", id)
	return id, nil
}
func (c *CompanyRepositoryImpl) DeleteCompany(ctx context.Context, companyId int64) error {
	query := `DELETE FROM companies WHERE company_id=$1`

	result, err := c.dbPoll.Exec(ctx, query, companyId)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		c.log.Error("Failed to delete company", "error", err)
		return dbErr
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("company %d does not exist", companyId)
	}
	c.log.Debug("Company deleted successfully", "id", companyId)
	return nil
}
