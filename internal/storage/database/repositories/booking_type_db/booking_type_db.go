package booking_type_db

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

var ErrBookingTypeNotFound = errors.New("Тип бронирования не найден ")

type BookingTypeRepository interface {
	CreateBookingType(ctx context.Context, dto services_models.CreateBookingTypeDTO) (int64, error)
	GetBookingType(ctx context.Context, BookingTypeId int64) (BookingTypeInfo, error)
	GetBookingTypeList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam, companyInfoDto services_models.CompanyInfo) (BookingTypeListResult, error)
	UpdateBookingType(ctx context.Context, id int64, name, description string) error
	DeleteBookingType(ctx context.Context, id int64) error
}
type BookingTypeInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type BookingTypeListResult struct {
	BookingTypes []BookingTypeInfo
	Total        int64
}

type BookingTypeRepositoryImpl struct {
	dbPoll *pgxpool.Pool
	log    *slog.Logger
}

func NewBookingTypeRepository(db *pgxpool.Pool, log *slog.Logger) *BookingTypeRepositoryImpl {
	return &BookingTypeRepositoryImpl{
		dbPoll: db,
		log:    log,
	}
}

func (bt *BookingTypeRepositoryImpl) CreateBookingType(ctx context.Context, dto services_models.CreateBookingTypeDTO) (int64, error) {
	query := `INSERT INTO booking_types (name, description, company_id) VALUES ($1, $2, $3) RETURNING id`

	var id int64
	err := bt.dbPoll.QueryRow(ctx, query, dto.Name, dto.Description, dto.CompanyId).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		bt.log.Error("Failed to create booking type", "error", err)
		return 0, database.PsqlErrorHandler(dbErr)
	}
	return id, nil
}

func (bt *BookingTypeRepositoryImpl) GetBookingType(ctx context.Context, BookingTypeId int64) (BookingTypeInfo, error) {
	query := `SELECT name, description FROM booking_types WHERE id = $1`

	var bookingType BookingTypeInfo
	err := bt.dbPoll.QueryRow(ctx, query, BookingTypeId).Scan(
		&bookingType.Name,
		&bookingType.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return BookingTypeInfo{}, ErrBookingTypeNotFound
	}
	if err != nil {
		if ctxErr := database.DbCtxError(ctx, err, bt.log); ctxErr != nil {
			return BookingTypeInfo{}, ctxErr
		}
		dbErr := database.PsqlErrorHandler(err)
		return BookingTypeInfo{}, dbErr
	}
	return bookingType, nil
}

func (bt *BookingTypeRepositoryImpl) GetBookingTypeList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam, companyInfoDto services_models.CompanyInfo) (BookingTypeListResult, error) {
	// Базовый SQL-запрос
	query := "SELECT id, name, description FROM booking_types WHERE company_id = $1"
	countQuery := "SELECT COUNT(*) FROM booking_types WHERE company_id = $1"
	searchQuery := " AND (name ILIKE $2 OR description ILIKE $2)"

	// Инициализация аргументов
	args := []interface{}{companyInfoDto.CompanyId}
	countArgs := []interface{}{companyInfoDto.CompanyId} // Начинаем с company_id для countQuery

	// Фильтрация по search
	if search != "" {
		query += searchQuery
		countQuery += searchQuery
		args = append(args, "%"+search+"%")
		countArgs = append(countArgs, "%"+search+"%") // Исправлено: добавляем в countArgs отдельно
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
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Подсчёт total
	var total int64
	err := bt.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total) // Используем countArgs
	if err != nil {
		bt.log.Error("Failed to count booking types", slog.Any("error", err))
		return BookingTypeListResult{}, fmt.Errorf("failed to count booking types: %w", err)
	}

	// Получение данных
	rows, err := bt.dbPoll.Query(ctx, query, args...)
	if err != nil {
		bt.log.Error("Failed to query booking types", slog.Any("error", err))
		return BookingTypeListResult{}, fmt.Errorf("failed to query booking types: %w", err)
	}
	defer rows.Close()

	var BookingTypes []BookingTypeInfo
	for rows.Next() {
		var BookingType BookingTypeInfo
		if err := rows.Scan(&BookingType.ID, &BookingType.Name, &BookingType.Description); err != nil {
			bt.log.Error("Error scanning booking type row", slog.Any("error", err))
			return BookingTypeListResult{}, fmt.Errorf("error scanning booking type row: %w", err)
		}
		BookingTypes = append(BookingTypes, BookingType)
	}
	if err = rows.Err(); err != nil {
		bt.log.Error("Error reading rows", slog.Any("error", err))
		return BookingTypeListResult{}, fmt.Errorf("error reading rows: %w", err)
	}

	return BookingTypeListResult{
		BookingTypes: BookingTypes,
		Total:        total,
	}, nil
}

func (bt *BookingTypeRepositoryImpl) UpdateBookingType(ctx context.Context, id int64, name, description string) error {
	query := `UPDATE booking_types SET name = $1, description = $2 WHERE id = $3`

	result, err := bt.dbPoll.Exec(ctx, query, name, description, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBookingTypeNotFound
		}
		dbErr := database.PsqlErrorHandler(err)
		bt.log.Error("Failed to update user in db", slog.String("error", err.Error()))
		return dbErr
	}
	if result.RowsAffected() == 0 {
		return ErrBookingTypeNotFound
	}
	bt.log.Debug("User updated successfully", "id", id)
	return nil
}

func (bt *BookingTypeRepositoryImpl) DeleteBookingType(ctx context.Context, id int64) error {
	query := `DELETE FROM booking_types WHERE id = $1`
	result, err := bt.dbPoll.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBookingTypeNotFound
		}
		dbErr := database.PsqlErrorHandler(err)
		bt.log.Error("Failed to delete user in db", slog.String("error", err.Error()))
		return dbErr
	}
	if result.RowsAffected() == 0 {
		return ErrBookingTypeNotFound
	}
	bt.log.Debug("Booking type deleted successfully", "id", id)
	return nil
}
