package jwt

import (
	"context"
	"fmt"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	Validate(ctx context.Context, tokenString string) (*Claims, error)
}

type hs256Validator struct {
	cfg JWTConfig
}

func NewValidator(cfg JWTConfig) Validator {
	return &hs256Validator{
		cfg: cfg,
	}
}

func (v *hs256Validator) Validate(
	ctx context.Context,
	tokenString string,
) (*Claims, error) {

	claims := &Claims{}

	token, err := jwtlib.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwtlib.Token) (interface{}, error) {
			// 🔐 Bắt buộc đúng thuật toán
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(v.cfg.Secret), nil
		},
		jwtlib.WithIssuer(v.cfg.Issuer), // ✅ check iss
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Name}),
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
