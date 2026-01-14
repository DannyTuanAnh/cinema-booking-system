package auth

import (
	"net/http"

	"cinema.com/demo/bff/clients/auth"
	"cinema.com/demo/bff/dto"
	"cinema.com/demo/bff/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authClient auth_clients.AuthClient
}

func NewAuthController(authClient auth_clients.AuthClient) *AuthController {
	return &AuthController{
		authClient: authClient,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var param dto.LoginDTO

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	resp, err := ac.authClient.Login(ctx, auth_clients.LoginRequest{
		Email:    param.Email,
		Password: param.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": resp})

}

func (ac *AuthController) Register(ctx *gin.Context) {
	var param dto.RegisterDTO

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}
	err := ac.authClient.Register(ctx, auth_clients.RegisterRequest{
		FullName: param.FullName,
		Email:    param.Email,
		Password: param.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}
