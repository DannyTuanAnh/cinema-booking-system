package movie_service

import (
	"context"

	"github.com/DannyTuanAnh/cinema-booking-system/internal/model"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/repository"
)

type MovieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository) *MovieService {
	return &MovieService{movieRepo: movieRepo}
}

func (s *MovieService) GetMovies(ctx context.Context) ([]model.Movie, error) {
	return s.movieRepo.GetMovies(ctx)
}
