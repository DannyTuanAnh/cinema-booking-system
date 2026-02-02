package routes

import (
	"os"

	"cinema.com/demo/bff/clients/movie"
	"cinema.com/demo/bff/controllers/movie"
	"github.com/gin-gonic/gin"
)

func InitMovieRoutes(r *gin.RouterGroup) {
	movieGroup := r.Group("/movies")
	RegisterMovieRoutes(movieGroup)
}

func RegisterMovieRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "http://" + addr + "/api"

	movieClient := movie_clients.NewMovieHTTPClient(path)

	movieController := movie.NewMovieController(movieClient)

	RegisterGetMovieRoute(r, movieController)
}

func RegisterGetMovieRoute(r *gin.RouterGroup, m *movie.MovieController) {
	r.GET("", m.GetMovie)
}
