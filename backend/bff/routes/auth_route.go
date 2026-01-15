package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/auth"
	"cinema.com/demo/bff/controllers/auth"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup, db *sql.DB) {
	authGroup := r.Group("/auth")
	RegisterAuthRoutes(authGroup, db)
}

func RegisterAuthRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	authClient := auth_clients.NewAuthHTTPClient(path)

	authController := auth.NewAuthController(authClient)

	RegisterSignUpRoute(r, authController, db)
	RegisterLoginRoute(r, authController, db)
}

func RegisterLoginRoute(r *gin.RouterGroup, ac *auth.AuthController, db *sql.DB) {
	r.POST("/login", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), ac.Login)
}

func RegisterSignUpRoute(r *gin.RouterGroup, ac *auth.AuthController, db *sql.DB) {
	r.POST("/register", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), ac.Register)
}
