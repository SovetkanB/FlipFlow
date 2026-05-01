package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
)

type contextKey string

const claimsKey contextKey = "claims"

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(w)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.jwtManager.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) *auth.Claims {
	v := ctx.Value(claimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*auth.Claims)
	return claims
}
