package booking_entity_db

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strings"
)

var ErrBookingEntityNotFound = errors.New("Объект бронирования не найден ")

type BookingEntityRepository interface {
	CreateBookingEntity(ctx context.Context, bookingTypeId int64, name, description, status string, ParentId int64) (int64, error)
	GetBookingEntity(ctx context.Context, BookingEntityId int64) (BookingEntityInfo, error)
	GetBookingEntitiesList(ctx context.Context, search string, limit, offset int, sort string) (BookingEntityListResult, error)
	UpdateBookingEntity(ctx context.Context, id int64, bookingTypeId int64, name, description, status string, ParentId int64) error
	DeleteBookingEntity(ctx context.Context, id int64) error
}
type BookingEntityInfo struct {
	ID            int64  `json:"id"`
	BookingTypeID int64  `json:"booking_type_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ParentID      int64  `json:"parent_id,omitempty"`
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

func (be *BookingEntityRepositoryImpl) CreateBookingEntity(ctx context.Context, bookingTypeId int64, name, description, status string, ParentId int64) (int64, error) {
	query := `INSERT INTO booking_entities (booking_type_id, name, description, parent_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := be.dbPoll.QueryRow(ctx, query, bookingTypeId, name, description, ParentId).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		be.log.Error("Failed to create booking entity", "error", err)
		return 0, database.PsqlErrorHandler(dbErr)
	}
	return id, nil
}

func (be *BookingEntityRepositoryImpl) GetBookingEntity(ctx context.Context, BookingEntityId int64) (BookingEntityInfo, error) {
	query := `SELECT booking_type_id, name, description, status, parent_id FROM booking_entities WHERE id = $1`

	var bookingEntity BookingEntityInfo
	err := be.dbPoll.QueryRow(ctx, query, BookingEntityId).Scan(
		&bookingEntity.BookingTypeID,
		&bookingEntity.Name,
		&bookingEntity.Description,
		&bookingEntity.Status,
		&bookingEntity.ParentID)
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

func (be *BookingEntityRepositoryImpl) GetBookingEntitiesList(ctx context.Context, search string, limit, offset int, sort string) (BookingEntityListResult, error) {
	// Базовый SQL-запрос для пользователей
	query := "SELECT id, booking_type_id, name, description, status, parent_id FROM booking_entities"
	countQuery := "SELECT COUNT(*) FROM booking_entities"
	//TODO Сделать сортировку как отдельные query
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
	if sort != "" {
		parts := strings.Split(sort, ":")
		if len(parts) == 2 && (parts[1] == "asc" || parts[1] == "desc") {
			// Простая проверка допустимых полей
			switch parts[0] {
			case "id", "description", "name":
				query += fmt.Sprintf(" ORDER BY %s %s", parts[0], strings.ToUpper(parts[1]))
			default:
				be.log.Warn("Invalid sort field", slog.String("field", parts[0]))
				return BookingEntityListResult{}, fmt.Errorf("invalid sort field: %s", parts[0])
			}
		} else {
			be.log.Warn("Invalid sort format", slog.String("sort", sort))
			return BookingEntityListResult{}, fmt.Errorf("invalid sort format: %s", sort)
		}
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Подсчёт total
	var total int64
	err := be.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		be.log.Error("Failed to count users", slog.Any("error", err))
		return BookingEntityListResult{}, fmt.Errorf("failed to count users: %w", err)
	}

	// Получение сущностей бронирования
	rows, err := be.dbPoll.Query(ctx, query, args...)
	if err != nil {
		be.log.Error("Failed to query users", slog.Any("error", err))
		return BookingEntityListResult{}, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var BookingEntities []BookingEntityInfo
	for rows.Next() {
		var BookingEntity BookingEntityInfo
		if err = rows.Scan(&BookingEntity.ID, &BookingEntity.BookingTypeID, &BookingEntity.Name, &BookingEntity.Description, &BookingEntity.Status, &BookingEntity.ParentID); err != nil {
			be.log.Error("Error scanning user row", slog.Any("error", err))
			return BookingEntityListResult{}, fmt.Errorf("error scanning user row: %w", err)
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
