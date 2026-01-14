package seat_services

import (
	"context"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
)

type SeatService struct {
	seatRepo repository.SeatRepository
}

func NewSeatService(seatRepo repository.SeatRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

func (s *SeatService) GetSeatByShowID(ctx context.Context, show_id int) ([]model.Seat, error) {
	return s.seatRepo.GetSeatByShowID(ctx, show_id)
}
