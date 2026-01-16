package book_service

import (
	"context"

	"log"

	"cinema.com/demo/internal/domain"
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

func (s *BookService) BookSeats(ctx context.Context, userID int64, seats []int) error {

	log.Println("Starting booking process for user:", userID, "with seats:", seats)
	tx, err := s.bookRepo.BeginTransaction(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	defer tx.Rollback()

	// set lock timeout
	err = s.bookRepo.SetTimeoutTx(ctx, tx, "3s")
	if err != nil {
		log.Println("Failed to set lock timeout:", err)
		return err
	}

	// lock the seats
	err = s.seatRepo.LockSeats(ctx, tx, seats)
	if err != nil {
		log.Println("Failed to lock seats:", err)
		return err
	}

	// check if all seats are still available
	count, err := s.seatRepo.CountSeatsForUpdate(ctx, tx, seats)
	if err != nil {
		log.Println("Failed to count seats for update:", err)
		return err
	}

	if count != len(seats) {
		return domain.ErrSeatAlreadyBooked
	}

	// book the seats
	err = s.seatRepo.BookSeats(ctx, tx, seats)
	if err != nil {
		log.Println("Failed to book seats:", err)
		return err
	}

	err = s.bookRepo.CreateBooking(ctx, tx, userID, seats)
	if err != nil {
		log.Println("Failed to create booking:", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}
