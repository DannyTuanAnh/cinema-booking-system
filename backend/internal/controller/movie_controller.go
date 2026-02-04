package controller

import (
	"net/http"

	movie_service "github.com/DannyTuanAnh/cinema-booking-system/internal/service/movie"
	"github.com/gin-gonic/gin"
)

type MovieController struct {
	movieService *movie_service.MovieService
}

func NewMovieController(movieService *movie_service.MovieService) *MovieController {
	return &MovieController{movieService: movieService}
}

func (c *MovieController) GetMovies(ctx *gin.Context) {
	movies, err := c.movieService.GetMovies(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}
