package auth

import "context"

type contextKey string

const ClaimsKey contextKey = "claims"

func ClaimsFromContext(ctx context.Context) *Claims {
	v := ctx.Value(ClaimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*Claims)
	return claims
}
