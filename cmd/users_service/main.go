package main

import (
	"context"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/config"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/models/booking_type/delete_booking_type"
	create_bookingType "github.com/ShlykovPavel/booker_microservice/internal/server/booking_type_handlers/create"
	"github.com/ShlykovPavel/booker_microservice/internal/server/booking_type_handlers/get_booking_types_list_handler"
	get_bookingType_by_id_handler "github.com/ShlykovPavel/booker_microservice/internal/server/booking_type_handlers/get_by_id"
	"github.com/ShlykovPavel/booker_microservice/internal/server/booking_type_handlers/update_booking_type"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database"
	"github.com/ShlykovPavel/booker_microservice/internal/storage/database/repositories/booking_type"
	users "github.com/ShlykovPavel/booker_microservice/user_service/server/users/create"
	users_delete "github.com/ShlykovPavel/booker_microservice/user_service/server/users/delete"
	"github.com/ShlykovPavel/booker_microservice/user_service/server/users/get_user"
	"github.com/ShlykovPavel/booker_microservice/user_service/server/users/get_user/get_user_list"
	"github.com/ShlykovPavel/booker_microservice/user_service/server/users/update_user"
	"github.com/ShlykovPavel/booker_microservice/user_service/storage/repositories/users_db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Info("Starting application")
	logger.Debug("Debug messages enabled")
	dbConfig := database.DbConfig{
		DbName:     cfg.DbName,
		DbUser:     cfg.DbUser,
		DbPassword: cfg.DbPassword,
		DbHost:     cfg.DbHost,
		DbPort:     cfg.DbPort,
	}

	poll, err := database.CreatePool(context.Background(), &dbConfig, logger)

	userRepository := users_db.NewUsersDB(poll, logger)
	bookerTypeRepository := booking_type.NewBookingTypeRepository(poll, logger)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/register", users.CreateUser(logger, userRepository, cfg.ServerTimeout))
	router.Get("/users/{id}", get_user.GetUserById(logger, userRepository, cfg.ServerTimeout))
	router.Get("/users", get_user_list.GetUserList(logger, userRepository, cfg.ServerTimeout))
	router.Put("/users/{id}", update_user.UpdateUserHandler(logger, userRepository, cfg.ServerTimeout))
	router.Delete("/users/{id}", users_delete.DeleteUserHandler(logger, userRepository, cfg.ServerTimeout))

	router.Post("/bookingType/create", create_bookingType.CreateBookingTypeHandler(logger, bookerTypeRepository, cfg.ServerTimeout))
	router.Get("/bookingType/{id}", get_bookingType_by_id_handler.GetBookingTypeByIdHandler(logger, bookerTypeRepository, cfg.ServerTimeout))
	router.Get("/bookingType", get_booking_types_list_handler.GetBookingTypesListHandler(logger, bookerTypeRepository, cfg.ServerTimeout))
	router.Put("/bookingType/{id}", update_booking_type.UpdateBookingTypeHandler(logger, bookerTypeRepository, cfg.ServerTimeout))
	router.Delete("/bookingType/{id}", delete_booking_type.DeleteBookingTypeHandler(logger, bookerTypeRepository, cfg.ServerTimeout))

	logger.Info("Starting HTTP server", slog.String("adress", cfg.Address))
	// Run server
	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.ServerTimeout,
		WriteTimeout:      cfg.ServerTimeout,
		//IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Stopped HTTP server")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}
