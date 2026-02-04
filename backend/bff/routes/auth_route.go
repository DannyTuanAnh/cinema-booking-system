package routes

import (
	"os"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/auth"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/controllers/auth"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	RegisterAuthRoutes(authGroup)
}

func RegisterAuthRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	authClient := auth_clients.NewAuthHTTPClient(path)

	authController := auth.NewAuthController(authClient)

	RegisterSignUpRoute(r, authController)
	RegisterLoginRoute(r, authController)
	RegisterRefreshTokenRoute(r, authController)
}

func RegisterLoginRoute(r *gin.RouterGroup, ac *auth.AuthController) {
	r.POST("/login", ac.Login)
}

func RegisterSignUpRoute(r *gin.RouterGroup, ac *auth.AuthController) {
	r.POST("/register", ac.Register)
}

func RegisterRefreshTokenRoute(r *gin.RouterGroup, ac *auth.AuthController) {
	r.POST("/refresh", ac.RefreshToken)
}
