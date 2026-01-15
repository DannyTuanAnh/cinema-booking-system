package auth_clients

import (
	"context"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}

type AuthClient interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, []*http.Cookie, error)
	Register(ctx context.Context, req RegisterRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, []*http.Cookie, error)
}
