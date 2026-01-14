package book_service

import (
	"context"

	"cinema.com/demo/internal/repository"
	"log"
)

type BookService struct {
	bookRepo repository.BookRepository
	seatRepo repository.SeatRepository
}

func NewBookService(bookRepo repository.BookRepository, seatRepo repository.SeatRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		seatRepo: seatRepo,
	}
}

func (s *BookService) BookSeats(ctx context.Context, userID int64, seats []int) error {

	log.Println("Starting booking process for user:", userID, "with seats:", seats)
	tx, err := s.bookRepo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = s.seatRepo.BookSeats(ctx, tx, userID, seats)
	if err != nil {
		return err
	}

	err = s.bookRepo.CreateBooking(ctx, tx, userID, seats)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
