package jwt

import (
	"context"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	GenerateAccessToken(ctx context.Context, userID int64, email, role string) (string, error)
}

type jwtHS256 struct {
	cfg JWTConfig
}

func NewJWTGenerator(cfg JWTConfig) JWTGenerator {
	return &jwtHS256{
		cfg: cfg,
	}
}

func (j *jwtHS256) GenerateAccessToken(ctx context.Context, userID int64, email, role string) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: int(userID),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwtLib.RegisteredClaims{
			Issuer:    j.cfg.Issuer,
			Subject:   string(rune(userID)),
			IssuedAt:  jwtLib.NewNumericDate(now),
			ExpiresAt: jwtLib.NewNumericDate(now.Add(j.cfg.Expire)),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.Secret))
}
