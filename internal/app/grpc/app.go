package grpcApp

import (
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/grpc_server"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

// В данном файле мы делаем обвязку нашего gRPC сервера, что б не нагружать файл main.go
// Это в каком то смысле приложение нашего gRPC сервера

// App Структура grpc сервера
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	// Создаём сущность сервера
	grpcServer := grpc.NewServer()

	//Подключаем к этой сущности наш обработчику
	grpc_server.Register(grpcServer)
	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	a.log.With(slog.String("op", "grpcApp.Run"), slog.Int("port", a.port))

	a.log.Info("starting grpc server")

	// делаем tcp слушатель на котором на gRPC будет принимать запросы
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	a.log.Info("gRPC server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (a *App) Stop() error {
	a.log.Info("stopping grpc server")

	a.gRPCServer.GracefulStop()

	a.log.Info("stopped grpc server")
	return nil
}
