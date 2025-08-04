package jwt_tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

// TokenClaims определяет структуру для claims из JWT-токена
type TokenClaims struct {
	AccountId     int64  `json:"AccountId"`
	CompanyId     int64  `json:"CompanyId"`
	CompanyCode   string `json:"CompanyCode"`
	CompanyLocale string `json:"CompanyLocale"`
	CompanyName   string `json:"CompanyName"`
	Email         string `json:"Email"`
	IsBot         string `json:"IsBot"`
	Phone         string `json:"Phone"`
	Role          string `json:"Role"`
	jwt.RegisteredClaims
}

// UnmarshalJSON реализует кастомную десериализацию для TokenClaims
// вызывается автоматически при парсинге токена в json в функции jwt.ParseWithClaims,
// так как реализует интерфейс UnmarshalJSON который используется библиотекой
func (c *TokenClaims) UnmarshalJSON(data []byte) error {
	// Вспомогательная структура для парсинга строковых значений
	type Alias struct {
		AccountId     string `json:"AccountId"`
		CompanyId     string `json:"CompanyId"`
		CompanyCode   string `json:"CompanyCode"`
		CompanyLocale string `json:"CompanyLocale"`
		CompanyName   string `json:"CompanyName"`
		Email         string `json:"Email"`
		IsBot         string `json:"IsBot"`
		Phone         string `json:"Phone"`
		Role          string `json:"Role"`
		jwt.RegisteredClaims
	}

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Преобразуем строковые AccountId и CompanyId в int
	accountId, err := strconv.ParseInt(alias.AccountId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid AccountId format: %w", err)
	}
	if accountId <= 0 {
		return fmt.Errorf("invalid AccountId value: %w", err)
	}
	companyId, err := strconv.ParseInt(alias.CompanyId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid CompanyId format: %w", err)
	}
	if companyId <= 0 {
		return fmt.Errorf("invalid companyId value: %w", err)
	}

	// Заполняем поля основной структуры
	c.AccountId = accountId
	c.CompanyId = companyId
	c.CompanyCode = alias.CompanyCode
	c.CompanyLocale = alias.CompanyLocale
	c.CompanyName = alias.CompanyName
	c.Email = alias.Email
	c.IsBot = alias.IsBot
	c.Phone = alias.Phone
	c.Role = alias.Role
	c.RegisteredClaims = alias.RegisteredClaims

	return nil
}

// VerifyToken verifies a JWT token and returns its claims
func VerifyToken(tokenString string, secretKey string) (*TokenClaims, error) {
	jwtClaims := &TokenClaims{}
	// Парсим токен
	tokenWithClaims, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	// Проверяем, валиден ли токен
	if !tokenWithClaims.Valid {
		return nil, errors.New("invalid token")
	}
	return jwtClaims, nil
}
