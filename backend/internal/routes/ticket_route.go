package routes

import (
	"github.com/DannyTuanAnh/cinema-booking-system/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitTicketRoutes(r *gin.RouterGroup, ticket *controller.TicketController) {
	m := r.Group("/tickets")
	{
		m.GET("", ticket.GetTicketByUserID)
	}
}
