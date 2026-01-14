package controller

import (
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

	token, userID, email, name, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, model.LoginResponse{
		AccessToken: token,
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
