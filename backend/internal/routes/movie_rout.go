package routes

import (
	"github.com/DannyTuanAnh/cinema-booking-system/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitMovieRoutes(r *gin.RouterGroup, movie *controller.MovieController) {
	m := r.Group("/movies")
	{
		m.GET("", movie.GetMovies)
	}
}
