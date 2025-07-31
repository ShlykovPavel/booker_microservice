package booking_db

import (
	"context"
	"errors"
	"fmt"
	create_booking_dto "github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking/create_booking"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/query_params"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strings"
	"time"
)

var ErrStartTimeAfterEndTime = errors.New("Start time after end time")
var ErrBookingNotFound = errors.New("Booking not found")

type BookingRepository interface {
	CreateBooking(ctx context.Context, dto create_booking_dto.BookingRequest) (int64, error)
	CheckBookingAvailability(ctx context.Context, bookingEntityId int64, startTime time.Time, endTime time.Time, excludeBookingId ...int64) (bool, error)
	GetBookingsByTime(ctx context.Context, startTime time.Time, endTime time.Time, queryParams query_params.ListQueryParams) ([]BookingInfo, error)
	GetBookingsByUserId(ctx context.Context, userId int64, queryParams query_params.ListQueryParams) (BookingList, error)
	GetBookingsByBookingEntity(ctx context.Context, BookingEntityId int64, queryParams query_params.ListQueryParams) (BookingList, error)
	GetBookingById(ctx context.Context, id int64) (BookingInfo, error)
	UpdateBooking(ctx context.Context, bookingInfo BookingInfo, bookingId int64) error
	DeleteBooking(ctx context.Context, bookingId int64) error
}
type BookingInfo struct {
	Id              int64
	UserId          int64
	BookingEntityId int64
	StartTime       time.Time
	EndTime         time.Time
	Status          string
}

type BookingList struct {
	Bookings []BookingInfo
	Total    int64
}

type BookingRepositoryImpl struct {
	dbPoll *pgxpool.Pool
	log    *slog.Logger
}

func NewBookingRepository(db *pgxpool.Pool, log *slog.Logger) *BookingRepositoryImpl {
	return &BookingRepositoryImpl{
		dbPoll: db,
		log:    log,
	}
}

func (b *BookingRepositoryImpl) CreateBooking(ctx context.Context, dto create_booking_dto.BookingRequest) (int64, error) {
	//Конвертация времени в UTC (если пришло не в UTC)
	startTime := dto.StartTime.UTC()
	endTime := dto.EndTime.UTC()

	if startTime.After(endTime) || startTime.Equal(endTime) {
		return 0, ErrStartTimeAfterEndTime
	}
	if dto.Status == "" {
		dto.Status = "pending"
	}
	query := `INSERT INTO bookings (user_id, booking_entity_id, start_time, end_time, status, company_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int64
	b.log.Debug("create booking sql request", "query", query)
	err := b.dbPoll.QueryRow(ctx, query, dto.UserId, dto.BookingEntityId, startTime, endTime, dto.Status, dto.CompanyId).Scan(&id)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		b.log.Error("Failed to create booking entity", "error", dbErr)
		return 0, dbErr
	}
	return id, nil
}

func (b *BookingRepositoryImpl) CheckBookingAvailability(ctx context.Context, bookingEntityId int64, startTime time.Time, endTime time.Time, excludeBookingId ...int64) (bool, error) {
	// Конвертируем время в UTC
	//startTime = startTime.UTC()
	//endTime = endTime.UTC()

	fmt.Println("CheckBookingAvailability: startTime", startTime, "endTime", endTime)

	if startTime.After(endTime) || startTime.Equal(endTime) {
		return false, ErrStartTimeAfterEndTime
	}

	query := `
        SELECT COUNT(*) 
        FROM bookings 
        WHERE booking_entity_id = $1 
        AND status != 'cancelled'
        AND (start_time, end_time) OVERLAPS ($2, $3)
    `
	args := []interface{}{bookingEntityId, startTime, endTime}

	// Если передан excludeBookingId, исключаем эту бронь из проверки
	if len(excludeBookingId) > 0 {
		query += ` AND id != $4`
		args = append(args, excludeBookingId[0])
	}

	var count int64
	b.log.Debug("check availability sql request", "query", query, "booking_entity_id", bookingEntityId, "start_time", startTime, "end_time", endTime)
	err := b.dbPoll.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		dbErr := database.PsqlErrorHandler(err)
		b.log.Error("Failed to check availability", "error", dbErr)
		return false, dbErr
	}

	// Если count == 0, интервал свободен
	return count == 0, nil
}

func (b *BookingRepositoryImpl) GetBookingsByTime(ctx context.Context, startTime time.Time, endTime time.Time, queryParams query_params.ListQueryParams) ([]BookingInfo, error) {
	//startTime = startTime.UTC()
	//endTime = endTime.UTC()

	if startTime.After(endTime) || startTime.Equal(endTime) {
		return nil, ErrStartTimeAfterEndTime
	}
	query := `SELECT id, user_id, booking_entity_id, start_time, end_time, status 
        FROM bookings 
        WHERE (start_time, end_time) OVERLAPS ($1, $2)`

	// Сортировка
	var orderBy []string
	if len(queryParams.SortParams) > 0 {
		for _, sortParam := range queryParams.SortParams {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sortParam.Field, strings.ToUpper(sortParam.Order)))
		}
		query += " ORDER BY " + strings.Join(orderBy, ", ")

	} else {
		// Дефолтная сортировка
		query += " ORDER BY start_time ASC"
	}

	b.log.Debug("get booking sql request", "query", query)
	rows, err := b.dbPoll.Query(ctx, query, startTime.UTC(), endTime.UTC())
	defer rows.Close()
	if err != nil {
		b.log.Error("Failed to get bookings by time", "error", err)
		return nil, database.PsqlErrorHandler(err)
	}

	var bookings []BookingInfo
	for rows.Next() {
		var bookingInfo BookingInfo
		if err = rows.Scan(&bookingInfo.Id, &bookingInfo.UserId, &bookingInfo.BookingEntityId, &bookingInfo.StartTime, &bookingInfo.EndTime, &bookingInfo.Status); err != nil {
			b.log.Error("Error scanning booking row", slog.Any("error", err))
			return nil, fmt.Errorf("error scanning booking row: %w", err)
		}
		bookings = append(bookings, bookingInfo)
	}
	return bookings, nil
}

func (b *BookingRepositoryImpl) GetBookingsByUserId(ctx context.Context, userId int64, queryParams query_params.ListQueryParams) (BookingList, error) {
	query := `SELECT id, user_id, booking_entity_id, start_time, end_time, status FROM bookings WHERE user_id = $1`
	countQuery := "SELECT COUNT(*) FROM bookings WHERE user_id = $1"
	args := []interface{}{userId}
	countArgs := []interface{}{userId}

	// Сортировка
	var orderBy []string
	if len(queryParams.SortParams) > 0 {
		for _, sortParam := range queryParams.SortParams {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sortParam.Field, strings.ToUpper(sortParam.Order)))
		}
		query += " ORDER BY " + strings.Join(orderBy, ", ")

	} else {
		// Дефолтная сортировка
		query += " ORDER BY id ASC"
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, queryParams.Limit, queryParams.Offset)

	// Подсчёт total
	var total int64
	err := b.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		b.log.Error("Failed to count users", slog.Any("error", err))
		return BookingList{}, fmt.Errorf("failed to count users: %w", err)
	}

	// Получение бронирований
	b.log.Debug("GetBookingsByUserId sql request", "query", query)
	rows, err := b.dbPoll.Query(ctx, query, args...)
	if err != nil {
		b.log.Error("Failed to query users", slog.Any("error", err))
		return BookingList{}, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var bookingsList []BookingInfo
	for rows.Next() {
		var bookingInfo BookingInfo
		if err = rows.Scan(&bookingInfo.Id, &bookingInfo.UserId, &bookingInfo.BookingEntityId, &bookingInfo.StartTime, &bookingInfo.EndTime, &bookingInfo.Status); err != nil {
			b.log.Error("Error scanning booking row", slog.Any("error", err))
			return BookingList{}, fmt.Errorf("failed to scan booking row: %w", err)
		}
		bookingInfo.StartTime = bookingInfo.StartTime.UTC()
		bookingInfo.EndTime = bookingInfo.EndTime.UTC()
		bookingsList = append(bookingsList, bookingInfo)
	}
	return BookingList{Bookings: bookingsList, Total: total}, nil
}

func (b *BookingRepositoryImpl) GetBookingsByBookingEntity(ctx context.Context, BookingEntityId int64, queryParams query_params.ListQueryParams) (BookingList, error) {
	query := `SELECT id, user_id, booking_entity_id, start_time, end_time, status FROM bookings WHERE booking_entity_id = $1`
	countQuery := "SELECT COUNT(*) FROM bookings WHERE booking_entity_id = $1"
	args := []interface{}{BookingEntityId}
	countArgs := []interface{}{BookingEntityId}

	// Сортировка
	var orderBy []string
	if len(queryParams.SortParams) > 0 {
		for _, sortParam := range queryParams.SortParams {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sortParam.Field, strings.ToUpper(sortParam.Order)))
		}
		query += " ORDER BY " + strings.Join(orderBy, ", ")

	} else {
		// Дефолтная сортировка
		query += " ORDER BY id ASC"
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, queryParams.Limit, queryParams.Offset)

	// Подсчёт total
	var total int64
	err := b.dbPoll.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		b.log.Error("Failed to count users", slog.Any("error", err))
		return BookingList{}, fmt.Errorf("failed to count users: %w", err)
	}

	// Получение бронирований
	rows, err := b.dbPoll.Query(ctx, query, args...)
	if err != nil {
		b.log.Error("Failed to query users", slog.Any("error", err))
		return BookingList{}, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var bookingsList []BookingInfo
	for rows.Next() {
		var bookingInfo BookingInfo
		if err = rows.Scan(&bookingInfo.Id, &bookingInfo.UserId, &bookingInfo.BookingEntityId, &bookingInfo.StartTime, &bookingInfo.EndTime, &bookingInfo.Status); err != nil {
			b.log.Error("Error scanning booking row", slog.Any("error", err))
			return BookingList{}, fmt.Errorf("failed to scan booking row: %w", err)
		}
		bookingInfo.StartTime = bookingInfo.StartTime.UTC()
		bookingInfo.EndTime = bookingInfo.EndTime.UTC()
		bookingsList = append(bookingsList, bookingInfo)
	}
	return BookingList{Bookings: bookingsList, Total: total}, nil
}

func (b *BookingRepositoryImpl) GetBookingById(ctx context.Context, id int64) (BookingInfo, error) {
	query := `SELECT id, user_id, booking_entity_id, start_time, end_time, status FROM bookings WHERE id = $1`
	b.log.Debug("get booking by id sql request", "query", query)

	var bookingInfo BookingInfo
	err := b.dbPoll.QueryRow(ctx, query, id).Scan(&bookingInfo.Id, &bookingInfo.UserId, &bookingInfo.BookingEntityId, &bookingInfo.StartTime, &bookingInfo.EndTime, &bookingInfo.Status)
	if err != nil {
		b.log.Error("Failed to get booking by id", "bookingId", id, "error", err)
		return BookingInfo{}, database.PsqlErrorHandler(err)
	}
	return bookingInfo, nil

}

func (b *BookingRepositoryImpl) UpdateBooking(ctx context.Context, bookingInfo BookingInfo, bookingId int64) error {
	query := `UPDATE bookings SET user_id =$1, booking_entity_id = $2, start_time = $3, end_time = $4, status =$5 WHERE id = $6`

	b.log.Debug("Updating booking sql request", "query", query)
	result, err := b.dbPoll.Exec(ctx, query, bookingInfo.UserId, bookingInfo.BookingEntityId, bookingInfo.StartTime, bookingInfo.EndTime, bookingInfo.Status, bookingId)
	if err != nil {
		b.log.Error("Error editing booking in db", slog.Any("error", err))
		return database.PsqlErrorHandler(err)
	}
	if result.RowsAffected() == 0 {
		return ErrBookingNotFound
	}
	b.log.Debug("Update booking successful ", "id", bookingInfo.Id, "user_id", bookingInfo.UserId)
	return nil
}

func (b *BookingRepositoryImpl) DeleteBooking(ctx context.Context, bookingId int64) error {
	query := `DELETE FROM bookings WHERE id = $1`

	b.log.Debug("Deleting booking sql request", "query", query)
	result, err := b.dbPoll.Exec(ctx, query, bookingId)
	if err != nil {
		b.log.Error("Error deleting editing booking in db", slog.Any("error", err))
		return database.PsqlErrorHandler(err)
	}
	if result.RowsAffected() == 0 {
		return ErrBookingNotFound
	}
	b.log.Debug("Delete booking successful ", "id", bookingId)
	return nil
}
