package routes

import (
	"cinema.com/demo/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitSeatRoutes(r *gin.RouterGroup, seat *controller.SeatController) {
	m := r.Group("/seats")
	{
		m.GET("", seat.GetSeatByShowID)
	}
}
