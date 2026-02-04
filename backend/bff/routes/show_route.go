package routes

import (
	"os"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/show"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/controllers/show"
	"github.com/gin-gonic/gin"
)

func InitShowRoutes(r *gin.RouterGroup) {
	showGroup := r.Group("/shows")
	RegisterShowRoutes(showGroup)
}

func RegisterShowRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	showClient := show_clients.NewShowHTTPClient(path)

	showController := show.NewShowController(showClient)

	RegisterGetShowRoute(r, showController)
}

func RegisterGetShowRoute(r *gin.RouterGroup, s *show.ShowController) {
	r.GET("", s.GetShowByMovieID)
}
