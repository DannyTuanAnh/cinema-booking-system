package book_service

import (
	"context"
	"fmt"

	"log"

	"cinema.com/demo/internal/repository"
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

func (s *BookService) BookSeats(ctx context.Context, userID int, seats []int) error {

	log.Println("Starting booking process for user:", userID, "with seats:", seats)
	tx, err := s.bookRepo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// set lock timeout
	err = s.bookRepo.SetTimeoutTx(ctx, tx, "3s")
	if err != nil {
		return err
	}

	// lock the seats
	err = s.seatRepo.LockSeats(ctx, tx, seats)
	if err != nil {
		return err
	}

	// check if all seats are still available
	count, err := s.seatRepo.CountSeatsForUpdate(ctx, tx, seats)
	if err != nil {
		return err
	}

	if count != len(seats) {
		return fmt.Errorf("one or more seats are already booked")
	}

	// book the seats

	err = s.seatRepo.BookSeats(ctx, tx, seats)
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
