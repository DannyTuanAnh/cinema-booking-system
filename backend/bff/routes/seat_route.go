package routes

import (
	"os"

	"cinema.com/demo/bff/clients/seat"
	"cinema.com/demo/bff/controllers/seat"
	"github.com/gin-gonic/gin"
)

func InitSeatRoutes(r *gin.RouterGroup) {
	seatGroup := r.Group("/seats")
	RegisterSeatRoutes(seatGroup)
}

func RegisterSeatRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	seatClient := seat_clients.NewSeatHTTPClient(path)

	seatController := seat.NewSeatController(seatClient)

	RegisterGetSeatRoute(r, seatController)
}

func RegisterGetSeatRoute(r *gin.RouterGroup, s *seat.SeatController) {
	r.GET("", s.GetSeatByShowID)
}
