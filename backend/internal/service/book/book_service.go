package book_service

import (
	"context"

	"log"

	"cinema.com/demo/internal/repository"
)

type BookService struct {
	seatRepo repository.SeatRepository
}

func NewBookService(seatRepo repository.SeatRepository) *BookService {
	return &BookService{
		seatRepo: seatRepo,
	}
}

func (s *BookService) BookSeats(ctx context.Context, userID int64, seats []int) error {
	// book the seats
	err := s.seatRepo.BookSeats(ctx, userID, seats)
	if err != nil {
		log.Println("Failed to book seats:", err)
		return err
	}

	return nil
}
