package main

import (
	"log"
	"os"
	// "strconv"
	"time"

	"github.com/DannyTuanAnh/cinema-booking-system/infra/db"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/controller"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/repository"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/routes"
	auth_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/auth"
	book_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/book"
	movie_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/movie"
	seat_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/seat"
	show_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/show"
	ticket_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/ticket"
	jwt "github.com/DannyTuanAnh/cinema-booking-system/pkg/jwt_service"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	// err = godotenv.Load("../../../.env")
	// if err != nil {
	// 	log.Println("No .env file found")
	// 	panic(err)
	// }

	// Database connection
	dbConfig := db.DefaultConfig()
	database, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	// expireMinutes, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_MINUTES"))

	jwtCfg := jwt.JWTConfig{
		Issuer: os.Getenv("JWT_ISSUER"),

		AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),

		// AccessTokenExpire:  time.Duration(expireMinutes) * time.Minute,
		// RefreshTokenExpire: 7 * time.Duration(expireHours) * time.Hour,

		AccessTokenExpire:  15 * time.Second,
		RefreshTokenExpire: 1 * time.Minute,
	}
	jwtGen := jwt.NewJWTGenerator(jwtCfg)
	jwtValid := jwt.NewValidator(jwtCfg)

	userRepo := repository.NewUserRepository(database)
	movieRepo := repository.NewMovieRepository(database)
	showRepo := repository.NewShowRepository(database)
	seatRepo := repository.NewSeatRepository(database)
	ticketRepo := repository.NewTicketRepository(database)
	rtRepo := repository.NewRefreshTokenRepoPG(database)

	authService := auth_service.NewAuthService(rtRepo, userRepo, jwtGen, jwtValid)
	movieService := movie_service.NewMovieService(movieRepo)
	showService := show_service.NewShowService(showRepo)
	seatService := seat_service.NewSeatService(seatRepo)
	bookService := book_service.NewBookService(seatRepo)
	ticketService := ticket_service.NewTicketService(ticketRepo)

	authController := controller.NewAuthController(authService)
	movieController := controller.NewMovieController(movieService)
	showController := controller.NewShowController(showService)
	seatController := controller.NewSeatController(seatService)
	bookController := controller.NewBookController(bookService)
	ticketController := controller.NewTicketController(ticketService)

	api := r.Group("/api")
	routes.InitAuthRoutes(api, authController)
	routes.InitMovieRoutes(api, movieController)
	routes.InitShowRoutes(api, showController)
	routes.InitSeatRoutes(api, seatController)
	routes.InitBookRoutes(api, bookController)
	routes.InitTicketRoutes(api, ticketController)

	addr := os.Getenv("ADDR_SERVER")

	if ENV := os.Getenv("ENV"); ENV == "development" {
		err = r.RunTLS(addr, "../../certs/localhost+2.pem", "../../certs/localhost+2-key.pem")
	} else {
		err = r.Run(addr)
	}

	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
