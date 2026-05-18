package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
)

func JWT(svc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				response.Unauthorized(w)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := svc.ValidateToken(tokenStr)
			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), auth.ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
