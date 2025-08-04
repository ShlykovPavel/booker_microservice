package booking_entity_db

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/services/services_models"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strings"
)

var ErrBookingEntityNotFound = errors.New("Объект бронирования не найден ")

type BookingEntityRepository interface {
	CreateBookingEntity(ctx context.Context, dto services_models.CreateBookingEntityDto) (int64, error)
	GetBookingEntity(ctx context.Context, BookingEntityId int64) (BookingEntityInfo, error)
	GetBookingEntitiesList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam, companyInfoDto services_models.CompanyInfo) (BookingEntityListResult, error)
	UpdateBookingEntity(ctx context.Context, id int64, bookingTypeId int64, name, description, status string, ParentId int64) error
	DeleteBookingEntity(ctx context.Context, id int64) error
	GetBookingTypeEntities(ctx context.Context, bookingTypeId int64, companyInfoDto services_models.CompanyInfo) ([]BookingEntityInfo, error)
}
type BookingEntityInfo struct {
	ID            int64  `json:"id"`
	BookingTypeID int64  `json:"booking_type_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ParentID      int64  `json:"parent_id,omitempty"`
	CompanyId     int64
}
type BookingEntityListResult struct {
	BookingEntities []BookingEntityInfo
	Total           int64
}

type BookingEntityRepositoryImpl struct {
	dbPoll *pgxpool.Pool
	log    *slog.Logger
}

func NewBookingEntityRepository(db *pgxpool.Pool, log *slog.Logger) *BookingEntityRepositoryImpl {
	return &BookingEntityRepositoryImpl{
		dbPoll: db,
		log:    log,
	}
}

func (be *BookingEntityRepositoryImpl) CreateBookingEntity(ctx context.Context, dto services_models.CreateBookingEntityDto) (int64, error) {
	query := `INSERT INTO booking_entities (booking_type_id, name, description, parent_id, company_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := be.dbPoll.QueryRow(ctx, query, dto.BookingEntityInfo.BookingTypeID, dto.BookingEntityInfo.Name, dto.BookingEntityInfo.Description, dto.BookingEntityInfo.ParentID, dto.CompanyId).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		be.log.Error("Failed to create booking entity", "error", err)
		return 0, database.PsqlErrorHandler(dbErr)
	}
	return id, nil
}

func (be *BookingEntityRepositoryImpl) GetBookingEntity(ctx context.Context, BookingEntityId int64) (BookingEntityInfo, error) {
	query := `SELECT booking_type_id, name, description, status, parent_id, company_id FROM booking_entities WHERE id = $1`

	var bookingEntity BookingEntityInfo
	err := be.dbPoll.QueryRow(ctx, query, BookingEntityId).Scan(
		&bookingEntity.BookingTypeID,
		&bookingEntity.Name,
		&bookingEntity.Description,
		&bookingEntity.Status,
		&bookingEntity.ParentID,
		&bookingEntity.CompanyId)
	if errors.Is(err, pgx.ErrNoRows) {
		return BookingEntityInfo{}, ErrBookingEntityNotFound
	}
	if err != nil {
		if ctxErr := database.DbCtxError(ctx, err, be.log); ctxErr != nil {
			return BookingEntityInfo{}, ctxErr
		}
		dbErr := database.PsqlErrorHandler(err)
		return BookingEntityInfo{}, dbErr
	}
	return bookingEntity, nil
}

func (be *BookingEntityRepositoryImpl) GetBookingEntitiesList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam, companyInfoDto services_models.CompanyInfo) (BookingEntityListResult, error) {
	// Базовый SQL-запрос
	query := "SELECT id, booking_type_id, name, description, status, parent_id FROM booking_entities WHERE company_id = $1"
	countQuery := "SELECT COUNT(*) FROM booking_entities WHERE company_id = $1"
	searchQuery := " AND (name ILIKE $2 OR description ILIKE $2)"

	// Инициализация аргументов
	args := []interface{}{companyInfoDto.CompanyId}
	countArgs := []interface{}{companyInfoDto.CompanyId}

	// Фильтрация по search
	if search != "" {
		query += searchQuery
		countQuery += searchQuery
		args = append(args, "%"+search+"%")
		countArgs = append(countArgs, "%"+search+"%")
	}

	// Сортировка
	var orderBy []string
	if len(sortParams) > 0 {
		for _, sortParam := range sortParams {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sortParam.Field, strings.ToUpper(sortParam.Order)))
		}
		query += " ORDER BY " + strings.Join(orderBy, ", ")
	} else {
		query += " ORDER BY id ASC"
	}

	// Пагинация
	paramOffset := len(args) + 1
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramOffset, paramOffset+1)
	args = append(args, limit, offset)

	// Подсчёт total
	var total int64
	err := be.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		be.log.Error("Failed to count booking entities", slog.Any("error", err))
		return BookingEntityListResult{}, fmt.Errorf("failed to count booking entities: %w", err)
	}

	// Получение данных
	rows, err := be.dbPoll.Query(ctx, query, args...)
	if err != nil {
		be.log.Error("Failed to query booking entities", slog.Any("error", err))
		return BookingEntityListResult{}, fmt.Errorf("failed to query booking entities: %w", err)
	}
	defer rows.Close()

	var BookingEntities []BookingEntityInfo
	for rows.Next() {
		var BookingEntity BookingEntityInfo
		if err = rows.Scan(
			&BookingEntity.ID,
			&BookingEntity.BookingTypeID,
			&BookingEntity.Name,
			&BookingEntity.Description,
			&BookingEntity.Status,
			&BookingEntity.ParentID,
		); err != nil {
			be.log.Error("Error scanning booking entity row", slog.Any("error", err))
			return BookingEntityListResult{}, fmt.Errorf("error scanning booking entity row: %w", err)
		}
		BookingEntities = append(BookingEntities, BookingEntity)
	}
	if err = rows.Err(); err != nil {
		be.log.Error("Error reading rows", slog.Any("error", err))
		return BookingEntityListResult{}, fmt.Errorf("error reading rows: %w", err)
	}

	return BookingEntityListResult{
		BookingEntities: BookingEntities,
		Total:           total,
	}, nil
}

func (be *BookingEntityRepositoryImpl) UpdateBookingEntity(ctx context.Context, id int64, bookingTypeId int64, name, description, status string, ParentId int64) error {
	query := `UPDATE booking_entities SET booking_type_id = $1, name = $2, description = $3, status = $4, parent_id = $5 WHERE id = $6`

	result, err := be.dbPoll.Exec(ctx, query, bookingTypeId, name, description, status, ParentId, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBookingEntityNotFound
		}
		dbErr := database.PsqlErrorHandler(err)
		be.log.Error("Failed to update booking entity in db", slog.String("error", err.Error()))
		return dbErr
	}
	if result.RowsAffected() == 0 {
		return ErrBookingEntityNotFound
	}
	be.log.Debug("booking entity updated successfully", "id", id)
	return nil
}

func (be *BookingEntityRepositoryImpl) DeleteBookingEntity(ctx context.Context, id int64) error {
	query := `DELETE FROM booking_entities WHERE id = $1`
	result, err := be.dbPoll.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBookingEntityNotFound
		}
		dbErr := database.PsqlErrorHandler(err)
		be.log.Error("Failed to delete booking entity in db", slog.String("error", err.Error()))
		return dbErr
	}
	if result.RowsAffected() == 0 {
		return ErrBookingEntityNotFound
	}
	be.log.Debug("booking entity deleted successfully", "id", id)
	return nil
}

// TODO Сделать поинт получения всех сущностей бронирования по id типа бронирования
func (be *BookingEntityRepositoryImpl) GetBookingTypeEntities(ctx context.Context, bookingTypeId int64, companyInfoDto services_models.CompanyInfo) ([]BookingEntityInfo, error) {
	query := `SELECT id, name, description FROM booking_entities WHERE booking_type_id = $1 AND company_id = $2`

	results, err := be.dbPoll.Query(ctx, query, bookingTypeId, companyInfoDto.CompanyId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBookingEntityNotFound
		}
		dbErr := database.PsqlErrorHandler(err)
		be.log.Error("Failed to get booking entities in db", slog.String("error", err.Error()))
		return nil, dbErr
	}
	defer results.Close()
	var bookingEntities []BookingEntityInfo
	for results.Next() {
		var bookingEntity BookingEntityInfo
		if err = results.Scan(
			&bookingEntity.ID,
			&bookingEntity.Name,
			&bookingEntity.Description); err != nil {
			be.log.Error("Failed to scan booking entity row", slog.Any("error", err))
			return nil, ErrBookingEntityNotFound
		}
		bookingEntities = append(bookingEntities, bookingEntity)

	}

	if err = results.Err(); err != nil {
		be.log.Error("Error reading rows", slog.Any("error", err))
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return bookingEntities, nil
}
