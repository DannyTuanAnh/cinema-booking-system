package auth_service

import (
	"context"
	"errors"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
	"cinema.com/demo/pkg/jwt_service"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	jwt      jwt.JWTGenerator
}

func NewAuthService(userRepo repository.UserRepository, jwt jwt.JWTGenerator) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, int64, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", 0, "", "", errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, "", "", errors.New("invalid password")
	}

	token, err := s.jwt.GenerateAccessToken(ctx, user.ID, user.Email, "user")
	if err != nil {
		return "", 0, "", "", err
	}

	return token, user.ID, user.Email, user.FullName, nil
}

func (s *AuthService) Register(ctx context.Context, fullName, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
