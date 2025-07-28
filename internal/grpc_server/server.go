package grpc_server

import (
	"context"
	pb "github.com/ShlykovPavel/booking-auth-proto/gen/go/booking_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BookingServiceServer реализует gRPC-сервис бронирований
type BookingServiceServer struct {
	pb.UnimplementedBookingServiceServer // Встроенная заглушка (обязательно!)
}

// Register Регистрирует наш обработчик (используем функцию регистрации в сгенерированном коде)
func Register(gRPC *grpc.Server) {
	pb.RegisterBookingServiceServer(gRPC, &BookingServiceServer{})
}

// NewBookingService создаёт экземпляр сервера
func NewBookingService() *BookingServiceServer {
	return &BookingServiceServer{}
}

// CreateUser функция создания пользователя.
//
// Тут мы берём типы из наших сгенерированных файлов pb.CreateUserRequest и pb.CreateUserResponse
func (s *BookingServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	// Проверяем, что email не пустой (базовая валидация)
	if req.Email == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"Email is required",
		)
	}

	// Заглушка: просто возвращаем тестовые данные
	return &pb.CreateUserResponse{
		Id:     1,
		Status: "User received: " + req.Email,
	}, nil
}
