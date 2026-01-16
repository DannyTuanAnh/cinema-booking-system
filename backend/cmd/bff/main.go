package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"cinema.com/demo/bff/middleware"
	"cinema.com/demo/bff/routes"
	"cinema.com/demo/bff/utils"
	"cinema.com/demo/infra/db"
	pkgFile "cinema.com/demo/pkg/file"
	jwt "cinema.com/demo/pkg/jwt_service"
	key "cinema.com/demo/pkg/key"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database connection
	dbConfig := db.DefaultConfig()
	database, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Register validator for struct
	if err := utils.RegisterValidator(); err != nil {
		panic(err)
	}

	fileEnv := "../../../.env"

	pkgFile.LoadEnv(fileEnv)

	// Check if command-line arguments are provided for API key generation
	if len(os.Args) >= 4 {
		clientType := os.Args[1]

		maxReq, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid maxReq parameter: %v", err)
		}

		winSec, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("Invalid winSec parameter: %v", err)
		}

		// Generate the API key
		err = key.GenerateAPIKey(fileEnv, clientType, maxReq, winSec, database)
		if err != nil {
			log.Fatalf("Failed to generate API key: %v", err)
		}

		return // Close the application after generating the API key
	}

	// expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	// expireMinutes, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_MINUTES"))

	jwtCfg := jwt.JWTConfig{
		Issuer: os.Getenv("JWT_ISSUER"),

		AccessSecret: os.Getenv("JWT_ACCESS_SECRET"),

		// AccessTokenExpire:  time.Duration(expireMinutes) * time.Minute,
		// RefreshTokenExpire: 7 * time.Duration(expireHours) * time.Hour,

		AccessTokenExpire:  15 * time.Second,
		RefreshTokenExpire: 1 * time.Minute,
	}

	jwtValidator := jwt.NewValidator(jwtCfg)
	jwtMiddleware := middleware.NewJWTMiddleware(jwtValidator)

	// Chạy server bình thường
	r := gin.Default()

	// Enable CORS for all routes
	r.Use(middleware.CORSMiddleware())

	// Serve frontend static files
	r.Static("/public", "../../../frontend")
	r.StaticFile("", "../../../frontend/index.html")

	api := r.Group("/api")

	// public
	routes.InitAuthRoutes(api, database)
	routes.InitMovieRoutes(api, database)
	routes.InitShowRoutes(api, database)
	routes.InitSeatRoutes(api, database)

	// protected
	protected := api.Group("")
	protected.Use(jwtMiddleware.Handle())
	{
		routes.InitBookRoutes(protected, database)
		routes.InitTicketRoutes(protected, database)
	}

	addr := os.Getenv("ADDR_BFF_SERVER")

	err = r.RunTLS(addr, "../../certs/localhost+2.pem", "../../certs/localhost+2-key.pem")

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
