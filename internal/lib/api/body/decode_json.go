package body

import (
	"errors"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// Ошибки
var ErrDecodeJSON = errors.New("failed to decode JSON")

// DecodeAndValidateJson декодирует JSON и валидирует структуру
func DecodeAndValidateJson(r *http.Request, v interface{}) error {
	// Декодируем JSON
	if err := render.DecodeJSON(r.Body, v); err != nil {
		slog.Default().Error("DecodeAndValidateJson: error decoding body or validating", "error", err)
		return ErrDecodeJSON
	}

	// Валидируем структуру
	if err := validator.New().Struct(v); err != nil {
		return err // Возвращаем ошибку валидации напрямую
	}

	return nil
}
