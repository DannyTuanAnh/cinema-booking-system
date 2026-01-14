package show_services

import (
	"context"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
)

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}

func (s *ShowService) GetShowByMovieID(ctx context.Context, movie_id int) ([]model.Show, error) {
	return s.showRepo.GetShowByMovieID(ctx, movie_id)
}
