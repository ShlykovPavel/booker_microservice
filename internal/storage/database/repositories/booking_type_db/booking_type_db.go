package booking_type_db

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strings"
)

var ErrBookingTypeNotFound = errors.New("Тип бронирования не найден ")

type BookingTypeRepository interface {
	CreateBookingType(ctx context.Context, name, description string) (int64, error)
	GetBookingType(ctx context.Context, BookingTypeId int64) (BookingTypeInfo, error)
	GetBookingTypeList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam) (BookingTypeListResult, error)
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

func (bt *BookingTypeRepositoryImpl) CreateBookingType(ctx context.Context, name, description string) (int64, error) {
	query := `INSERT INTO booking_types (name, description) VALUES ($1, $2) RETURNING id`

	var id int64
	err := bt.dbPoll.QueryRow(ctx, query, name, description).Scan(&id)
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

func (bt *BookingTypeRepositoryImpl) GetBookingTypeList(ctx context.Context, search string, limit, offset int, sortParams []query_params.SortParam) (BookingTypeListResult, error) {
	// Базовый SQL-запрос для пользователей
	query := "SELECT id, name, description FROM booking_types"
	countQuery := "SELECT COUNT(*) FROM booking_types"
	searchQuery := " WHERE name ILIKE $1 OR description ILIKE $1"
	args := []interface{}{}
	countArgs := []interface{}{}

	// Фильтрация по search
	if search != "" {
		query += searchQuery
		countQuery += searchQuery
		args = append(args, "%"+search+"%")
		countArgs = append(countArgs, "%"+search+"%")
	}

	// Сортировка
	//TODO Изменить сортировку на отдельные query
	var orderBy []string
	if len(sortParams) > 0 {
		for _, sortParam := range sortParams {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sortParam.Field, strings.ToUpper(sortParam.Order)))
		}
		query += " ORDER BY " + strings.Join(orderBy, ", ")

	} else {
		// Дефолтная сортировка
		query += " ORDER BY id ASC"
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Подсчёт total
	var total int64
	err := bt.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		bt.log.Error("Failed to count users", slog.Any("error", err))
		return BookingTypeListResult{}, fmt.Errorf("failed to count users: %w", err)
	}

	// Получение пользователей
	rows, err := bt.dbPoll.Query(ctx, query, args...)
	if err != nil {
		bt.log.Error("Failed to query users", slog.Any("error", err))
		return BookingTypeListResult{}, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var BookingTypes []BookingTypeInfo
	for rows.Next() {
		var BookingType BookingTypeInfo
		if err := rows.Scan(&BookingType.ID, &BookingType.Name, &BookingType.Description); err != nil {
			bt.log.Error("Error scanning user row", slog.Any("error", err))
			return BookingTypeListResult{}, fmt.Errorf("error scanning user row: %w", err)
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
	bt.log.Debug("User deleted successfully", "id", id)
	return nil
}
