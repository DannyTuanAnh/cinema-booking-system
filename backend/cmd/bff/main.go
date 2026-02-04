package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	MiddlewareApiKey "github.com/DannyTuanAnh/cinema-booking-system/bff/middleware/api_key"
	MiddlewareCors "github.com/DannyTuanAnh/cinema-booking-system/bff/middleware/cors"
	MiddlewareJWT "github.com/DannyTuanAnh/cinema-booking-system/bff/middleware/jwt"
	MiddlewareRateLimit "github.com/DannyTuanAnh/cinema-booking-system/bff/middleware/rate_limit"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/routes"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/utils"
	"github.com/DannyTuanAnh/cinema-booking-system/infra/db"
	pkgFile "github.com/DannyTuanAnh/cinema-booking-system/pkg/file"
	jwt "github.com/DannyTuanAnh/cinema-booking-system/pkg/jwt_service"
	key "github.com/DannyTuanAnh/cinema-booking-system/pkg/key"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables from .env file
	fileEnv := "../../../.env"
	pkgFile.LoadEnv(fileEnv)

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
	jwtMiddleware := MiddlewareJWT.NewJWTMiddleware(jwtValidator)

	// Chạy server bình thường
	r := gin.Default()

	// Enable CORS for all routes
	r.Use(MiddlewareCors.CORSMiddleware())

	// Khởi tạo Redis client cho Rate Limiting
	rds := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS")})
	redisLimiter := MiddlewareRateLimit.NewRedisRateLimiter(rds, 60, 100) // 100 requests per 60 seconds
	normalLimiter := MiddlewareRateLimit.NewNormalRateLimiter()

	// Khởi tạo RedisHealth để theo dõi trạng thái của Redis
	redisHealth := MiddlewareRateLimit.NewRedisHealth()

	// Bắt đầu kiểm tra trạng thái của Redis mỗi 3 giây
	utils.StartRedisHealthChecker(context.Background(), rds, redisHealth, 3*time.Second)

	api := r.Group("/api")
	api.Use(jwtMiddleware.ExtractUser(), MiddlewareRateLimit.RateLimitMiddleware(redisLimiter, normalLimiter, redisHealth), MiddlewareApiKey.ApiKeyMiddleware(database))
	{

		// public
		routes.InitAuthRoutes(api)
		routes.InitMovieRoutes(api)
		routes.InitShowRoutes(api)
		routes.InitSeatRoutes(api)

		// protected
		protected := api.Group("")
		protected.Use(jwtMiddleware.Handle())
		{
			routes.InitBookRoutes(protected)
			routes.InitTicketRoutes(protected)
		}
	}

	addr := os.Getenv("ADDR_BFF_SERVER")

	if ENV := os.Getenv("ENV"); ENV == "development" {
		err = r.RunTLS(addr, "../../certs/localhost+2.pem", "../../certs/localhost+2-key.pem")
	} else {
		err = r.Run(addr)
	}

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
