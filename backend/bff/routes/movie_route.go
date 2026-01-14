package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/movie"
	"cinema.com/demo/bff/controllers/movie"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitMovieRoutes(r *gin.RouterGroup, db *sql.DB) {
	movieGroup := r.Group("/movies")
	RegisterMovieRoutes(movieGroup, db)
}

func RegisterMovieRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "http://" + addr + "/api"

	movieClient := movie_clients.NewMovieHTTPClient(path)

	movieController := movie.NewMovieController(movieClient)

	RegisterGetMovieRoute(r, movieController, db)
}

func RegisterGetMovieRoute(r *gin.RouterGroup, m *movie.MovieController, db *sql.DB) {
	r.GET("", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), m.GetMovie)
}
