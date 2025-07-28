package app

import (
	grpcApp2 "github.com/ShlykovPavel/booker_microservice/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcApp2.App
}

func NewApp(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcApp2.NewApp(logger, grpcPort)
	return &App{GRPCServer: grpcApp}
}
