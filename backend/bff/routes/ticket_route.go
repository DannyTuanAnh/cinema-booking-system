package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/ticket"
	"cinema.com/demo/bff/controllers/ticket"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitTicketRoutes(r *gin.RouterGroup, db *sql.DB) {
	ticketGroup := r.Group("/tickets")
	RegisterTicketRoutes(ticketGroup, db)
}

func RegisterTicketRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "http://" + addr + "/api"

	ticketClient := ticket_clients.NewTicketHTTPClient(path)

	ticketController := ticket.NewTicketController(ticketClient)

	RegisterGetTicketRoute(r, ticketController, db)
}

func RegisterGetTicketRoute(r *gin.RouterGroup, t *ticket.TicketController, db *sql.DB) {
	r.GET("", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), t.GetTicketByUserID)
}
