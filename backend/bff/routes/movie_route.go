package routes

import (
	"os"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/movie"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/controllers/movie"
	"github.com/gin-gonic/gin"
)

func InitMovieRoutes(r *gin.RouterGroup) {
	movieGroup := r.Group("/movies")
	RegisterMovieRoutes(movieGroup)
}

func RegisterMovieRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := addr + "/api"

	movieClient := movie_clients.NewMovieHTTPClient(path)

	movieController := movie.NewMovieController(movieClient)

	RegisterGetMovieRoute(r, movieController)
}

func RegisterGetMovieRoute(r *gin.RouterGroup, m *movie.MovieController) {
	r.GET("", m.GetMovie)
}
