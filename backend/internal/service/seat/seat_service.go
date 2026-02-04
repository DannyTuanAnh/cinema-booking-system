package seat_services

import (
	"context"

	"github.com/DannyTuanAnh/cinema-booking-system/internal/model"
	"github.com/DannyTuanAnh/cinema-booking-system/internal/repository"
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
