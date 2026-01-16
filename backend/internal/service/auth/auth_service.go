package auth_service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
	jwt "cinema.com/demo/pkg/jwt_service"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrReuseDetected     = errors.New("refresh token reuse detected")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrInvalidToken      = errors.New("invalid token")
)

type AuthService struct {
	rtRepo   repository.RefreshTokenRepo
	userRepo repository.UserRepository
	jwtGen   jwt.JWTGenerator
	jwtValid jwt.Validator
}

func NewAuthService(rtRepo repository.RefreshTokenRepo, userRepo repository.UserRepository, jwtGen jwt.JWTGenerator, jwtValid jwt.Validator) *AuthService {
	return &AuthService{
		rtRepo:   rtRepo,
		userRepo: userRepo,
		jwtGen:   jwtGen,
		jwtValid: jwtValid,
	}
}

func (s *AuthService) RefreshTokenMaxAge() time.Duration {
	return s.jwtGen.RefreshTokenMaxAge()
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, int64, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", 0, "", "", errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", 0, "", "", errors.New("invalid password")
	}

	accessToken, err := s.jwtGen.GenerateAccessToken(ctx, user.ID, user.Email, "user")
	if err != nil {
		return "", "", 0, "", "", err
	}

	refreshToken, err := s.jwtGen.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", 0, "", "", err
	}

	err = s.rtRepo.Save(ctx, &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: repository.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.jwtGen.RefreshTokenMaxAge()),
	})
	if err != nil {
		return "", "", 0, "", "", err
	}

	return accessToken, refreshToken, user.ID, user.Email, user.FullName, nil
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

func (s *AuthService) Refresh(ctx context.Context, oldToken string) (string, string, error) {

	// verify jwt old refresh token

	claims, err := s.jwtValid.ValidateRefresh(ctx, oldToken)
	if err != nil {
		//token giả
		return "", "", ErrInvalidCredential
	}

	userIDFromToken, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		//token giả
		return "", "", ErrInvalidToken
	}

	oldHash := repository.HashToken(oldToken)

	rt, err := s.rtRepo.FindByHash(ctx, oldHash)
	if err != nil {
		//token giả
		return "", "", ErrInvalidCredential
	}

	// Kiểm tra userID từ old refresh token và từ db có trùng không
	if userIDFromToken != rt.UserID {
		return "", "", ErrInvalidCredential
	}

	if rt.Revoked || time.Now().After(rt.ExpiresAt) {

		//reuse detected
		err := s.rtRepo.RevokeAllByUser(ctx, rt.UserID)
		if err != nil {
			return "", "", err
		}

		return "", "", ErrReuseDetected

	}

	// revoke old
	newRefresh, err := s.jwtGen.GenerateRefreshToken(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}

	newHash := repository.HashToken(newRefresh)

	err = s.rtRepo.Revoke(ctx, rt.ID, newHash)
	if err != nil {
		return "", "", err
	}

	// save new
	err = s.rtRepo.Save(ctx, &model.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(s.jwtGen.RefreshTokenMaxAge()),
	})
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}

	access, _ := s.jwtGen.GenerateAccessToken(
		ctx,
		user.ID,
		user.Email,
		"user",
	)

	return access, newRefresh, nil
}
