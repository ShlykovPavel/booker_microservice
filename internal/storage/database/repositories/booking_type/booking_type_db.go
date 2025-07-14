package booking_type

import (
	"context"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type BookingTypeRepository interface {
	CreateBookingType(ctx context.Context, name, description string) (int64, error)
}
type BookingTypeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
