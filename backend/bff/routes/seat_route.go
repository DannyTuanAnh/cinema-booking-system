package routes

import (
	"os"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/seat"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/controllers/seat"
	"github.com/gin-gonic/gin"
)

func InitSeatRoutes(r *gin.RouterGroup) {
	seatGroup := r.Group("/seats")
	RegisterSeatRoutes(seatGroup)
}

func RegisterSeatRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := addr + "/api"

	seatClient := seat_clients.NewSeatHTTPClient(path)

	seatController := seat.NewSeatController(seatClient)

	RegisterGetSeatRoute(r, seatController)
}

func RegisterGetSeatRoute(r *gin.RouterGroup, s *seat.SeatController) {
	r.GET("", s.GetSeatByShowID)
}
