package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cinema.com/demo/bff/middleware"
	"cinema.com/demo/bff/routes"
	"cinema.com/demo/bff/utils"
	"cinema.com/demo/infra/db"
	jwt "cinema.com/demo/pkg/jwt_service"

	"strings"

	repo "cinema.com/demo/bff/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Database connection
	dbConfig := db.DefaultConfig()
	database, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := utils.RegisterValidator(); err != nil {
		panic(err)
	}

	err = godotenv.Load("../../../.env")
	if err != nil {
		log.Println("No .env file found")
		panic(err)
	}

	// Kiểm tra nếu có arguments để tạo API key
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

		// Tạo API key
		apiRepo := repo.NewApiKeyRepo(database)

		exits, _ := apiRepo.ExistsByClient(clientType)
		if !exits {
			plaintext := utils.GenerateApiKey(clientType)
			hash := utils.HashApiKey(plaintext)

			if err := apiRepo.Insert(clientType, hash, maxReq, winSec); err != nil {
				log.Fatal(err)
			}
			if clientType == "web" {

				// Append to .env
				envLine := fmt.Sprintf("\nVITE_API_KEY=%s\n", plaintext)

				f, err := os.OpenFile("../../../frontend/.env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				if _, err := f.WriteString(envLine); err != nil {
					log.Fatal(err)
				}

				log.Printf("API key created successfully for client: %s", clientType)
				log.Printf("VITE_API_KEY=%s", plaintext)

			} else {

				// Append to .env
				envLine := fmt.Sprintf("\nAPI_KEY_%s=%s\n", strings.ToUpper(clientType), plaintext)

				f, err := os.OpenFile("../../../frontend/.env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				if _, err := f.WriteString(envLine); err != nil {
					log.Fatal(err)
				}

				log.Printf("API key created successfully for client: %s", clientType)
				log.Printf("API_KEY_%s=%s", strings.ToUpper(clientType), plaintext)
			}
		} else {
			log.Printf("API key already exists for client: %s", clientType)
		}

		return // Kết thúc chương trình sau khi tạo API key
	}

	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	jwtCfg := jwt.JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: os.Getenv("JWT_ISSUER"),
		Expire: time.Duration(expireHours) * time.Hour,
	}

	jwtValidator := jwt.NewValidator(jwtCfg)
	jwtMiddleware := middleware.NewJWTMiddleware(jwtValidator)

	// Chạy server bình thường
	r := gin.Default()

	// Enable CORS for all routes
	r.Use(middleware.CORSMiddleware())

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

	r.Run(addr)
}
