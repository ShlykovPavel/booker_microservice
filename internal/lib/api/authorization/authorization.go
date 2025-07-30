package authorization

import (
	"github.com/ShlykovPavel/booker_microservice/internal/lib/jwt_tokens"
)

// Authorization проверяет предоставленный токен и получает аргументы тела токена
func Authorization(tokenString string, secretKey string) (*jwt_tokens.TokenClaims, error) {
	jwtClaims, err := jwt_tokens.VerifyToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}
	return jwtClaims, nil
}
