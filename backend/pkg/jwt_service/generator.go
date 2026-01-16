package jwt

import (
	"context"
	"log"
	"strconv"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	GenerateAccessToken(ctx context.Context, userID int64, email, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID int64) (string, error)

	RefreshTokenMaxAge() time.Duration
}

type jwtHS256 struct {
	cfg JWTConfig
}

func NewJWTGenerator(cfg JWTConfig) JWTGenerator {
	return &jwtHS256{cfg: cfg}
}

func (j *jwtHS256) GenerateAccessToken(ctx context.Context, userID int64, email, role string) (string, error) {
	now := time.Now()

	claims := &AccessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwtLib.RegisteredClaims{
			Issuer:    j.cfg.Issuer,
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwtLib.NewNumericDate(now),
			ExpiresAt: jwtLib.NewNumericDate(now.Add(j.cfg.AccessTokenExpire)),
		},
	}

	log.Println("access token claims:", claims)

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	log.Println("generating access token with claims:", token)

	return token.SignedString([]byte(j.cfg.AccessSecret))
}

func (j *jwtHS256) GenerateRefreshToken(ctx context.Context, userID int64) (string, error) {
	now := time.Now()

	claims := jwtLib.RegisteredClaims{
		Issuer:    j.cfg.Issuer,
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwtLib.NewNumericDate(now),
		ExpiresAt: jwtLib.NewNumericDate(now.Add(j.cfg.RefreshTokenExpire)),
	}

	log.Println("refresh token claims:", claims)

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	log.Println("generating refresh token with claims:", token)

	return token.SignedString([]byte(j.cfg.RefreshSecret))
}

func (j *jwtHS256) RefreshTokenMaxAge() time.Duration {
	return j.cfg.RefreshTokenExpire
}
