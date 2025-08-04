package helpers

import (
	"context"
	resp "github.com/ShlykovPavel/booker_microservice/internal/lib/api/response"
	"github.com/ShlykovPavel/booker_microservice/internal/lib/jwt_tokens"
	"log/slog"
	"net/http"
)

func ExtractTokenClaims(ctx context.Context, log *slog.Logger, w http.ResponseWriter, r *http.Request) *jwt_tokens.TokenClaims {
	claimsRaw := ctx.Value("tokenClaims")
	if claimsRaw == nil {
		resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid auth token: Claims not found"))
		return nil
	}

	claims, ok := claimsRaw.(*jwt_tokens.TokenClaims)
	if !ok {
		log.Error("Invalid claims type in context")
		resp.RenderResponse(w, r, http.StatusBadRequest, resp.Error("Invalid auth token: Invalid claims type"))
		return nil
	}
	log.Debug("Получены token claims", slog.Any("claims", claims))
	return claims
}
