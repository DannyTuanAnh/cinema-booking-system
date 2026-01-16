package controller

import (
	"log"
	"net/http"

	"cinema.com/demo/internal/model"
	auth_service "cinema.com/demo/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *auth_service.AuthService
}

func NewAuthController(authService *auth_service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, userID, email, name, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
		MaxAge:   int(c.authService.RefreshTokenMaxAge().Seconds()),
	})

	ctx.JSON(http.StatusOK, model.LoginResponse{
		AccessToken: access,
		UserID:      userID,
		Email:       email,
		Name:        name,
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.authService.Register(
		ctx.Request.Context(),
		req.FullName,
		req.Email,
		req.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	rt, err := ctx.Cookie("refresh_token")

	log.Println("old refresh token cookie from client: ", rt)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	access, refresh, err := c.authService.Refresh(ctx, rt)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token reuse detected"})
		return
	}

	log.Println("new access token:", access)
	log.Println("new refresh token:", refresh)

	maxAge := int(c.authService.RefreshTokenMaxAge().Seconds())

	// set new refresh cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
		MaxAge:   maxAge,
	})

	ctx.JSON(http.StatusOK, gin.H{"access_token": access})
}
