package movie_service

import (
	"context"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
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
