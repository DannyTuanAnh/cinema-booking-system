package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/seat"
	"cinema.com/demo/bff/controllers/seat"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitSeatRoutes(r *gin.RouterGroup, db *sql.DB) {
	seatGroup := r.Group("/seats")
	RegisterSeatRoutes(seatGroup, db)
}

func RegisterSeatRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	seatClient := seat_clients.NewSeatHTTPClient(path)

	seatController := seat.NewSeatController(seatClient)

	RegisterGetSeatRoute(r, seatController, db)
}

func RegisterGetSeatRoute(r *gin.RouterGroup, s *seat.SeatController, db *sql.DB) {
	r.GET("", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), s.GetSeatByShowID)
}
