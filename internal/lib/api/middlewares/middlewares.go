package middlewares

import (
	"context"
	"fmt"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/api/authorization"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/jwt_tokens"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
)

// AuthMiddleware проверяет токен авторизации при выполнении запроса
//
// # При успехе передаёт обработку следующему хендлеру
//
// При ошибке возвращает статус код 401 и ошибку
func AuthMiddleware(secretKey string, log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "internal/lib/api/middlewares/middlewares.go/AuthMiddleware"
	log = log.With(slog.String("op", op))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				renderUnauthorized(w, r, log, "Authorization header is missing")
				return
			}
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				renderUnauthorized(w, r, log, "Authorization header is invalid")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := authorization.Authorization(tokenString, secretKey)
			if err != nil {
				renderUnauthorized(w, r, log, fmt.Sprintf("Authorization token is invalid: %v", err))
				return
			}
			log.Debug("Authorization token is valid", slog.Any("claims", claims))
			ctx := context.WithValue(r.Context(), "tokenClaims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}
}

func renderUnauthorized(w http.ResponseWriter, r *http.Request, log *slog.Logger, msg string) {
	log.Error(msg)
	render.Status(r, http.StatusUnauthorized)
	render.JSON(w, r, resp.Error(msg))
}

func AuthAdminMiddleware(secretKey string, log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "internal/lib/api/middlewares/middlewares.go/AuthAdminMiddleware"
	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		// Используем AuthMiddleware для проверки авторизации
		return AuthMiddleware(secretKey, log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем claims из контекста
			claims, ok := r.Context().Value("tokenClaims").(*jwt_tokens.TokenClaims)
			if !ok {
				log.Error("Failed to retrieve claims from context")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("Internal server error"))
				return
			}
			value := claims.Role

			// Проверяем, является ли пользователь администратором
			switch value {
			case "Administrator":
				log.Debug("User is authorized and has admin privileges")
				next.ServeHTTP(w, r)
			case "Owner":
				log.Debug("User is authorized and has admin privileges")
				next.ServeHTTP(w, r)
			case "Manager":
				log.Debug("User is authorized and has admin privileges")
				next.ServeHTTP(w, r)
			default:
				log.Debug("User is NOT authorized and has admin privileges")
				resp.RenderResponse(w, r, http.StatusForbidden,
					resp.Error("Forbidden"))
			}
		}))
	}
}
