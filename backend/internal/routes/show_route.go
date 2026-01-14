package routes

import (
	"cinema.com/demo/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitShowRoutes(r *gin.RouterGroup, show *controller.ShowController) {
	m := r.Group("/shows")
	{
		m.GET("", show.GetShowByMovieID)
	}
}
