package movie

import (
	"net/http"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/movie"
	"github.com/gin-gonic/gin"
)

type MovieController struct {
	movieClient movie_clients.MovieClient
}

func NewMovieController(movieClient movie_clients.MovieClient) *MovieController {
	return &MovieController{
		movieClient: movieClient,
	}
}

func (mc *MovieController) GetMovie(ctx *gin.Context) {
	resp, err := mc.movieClient.GetMovieDetails(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": resp})
}
