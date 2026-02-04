package show_services

import (
	"context"

	"github.com/DannyTuanAnh/cinema-booking-system/internal/model"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/repository"
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
