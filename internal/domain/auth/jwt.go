package auth

import (
	"fmt"
	"time"

	"github.com/SovetkanB/FlipFlow/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(jwtconfig config.JWTConfig, user *User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(jwtconfig.AccessTTL)

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtconfig.Secret))
}

func ValidateToken(jwtconfig config.JWTConfig, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtconfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claim, nil
}
