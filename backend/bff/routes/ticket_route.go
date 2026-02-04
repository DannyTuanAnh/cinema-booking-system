package routes

import (
	"os"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/ticket"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/controllers/ticket"
	"github.com/gin-gonic/gin"
)

func InitTicketRoutes(r *gin.RouterGroup) {
	ticketGroup := r.Group("/tickets")
	RegisterTicketRoutes(ticketGroup)
}

func RegisterTicketRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	ticketClient := ticket_clients.NewTicketHTTPClient(path)

	ticketController := ticket.NewTicketController(ticketClient)

	RegisterGetTicketRoute(r, ticketController)
}

func RegisterGetTicketRoute(r *gin.RouterGroup, t *ticket.TicketController) {
	r.GET("", t.GetTicketByUserID)
}
