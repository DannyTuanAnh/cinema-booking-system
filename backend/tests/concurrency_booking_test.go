package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"cinema.com/demo/internal/repository"
	book_service "cinema.com/demo/internal/service/book"
)

func setupDB(t *testing.T) *sql.DB {
	dsn := "postgres://postgres:1@localhost:5432/cinema?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func TestConcurrentSeatBooking(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	ctx := context.Background()

	seatRepo := repository.NewSeatRepository(db)
	bookRepo := repository.NewBookRepository(db)
	bookService := book_service.NewBookService(bookRepo, seatRepo)

	// seats := []int{7, 8, 9}
	userIDs := []int64{4, 5, 6}

	var wg sync.WaitGroup
	wg.Add(len(userIDs))

	results := make([]error, 3)

	start := make(chan struct{ userId int64 })

	for i, userID := range userIDs {
		go func(idx int, uid int64) {
			defer wg.Done()

			<-start

			err := bookService.BookSeats(ctx, uid, []int{7 + idx, 8 + idx, 9 + idx})
			if err != nil {
				results[idx] = fmt.Errorf("user %d, %w", uid, err)
			} else {
				results[idx] = nil
			}
		}(i, userID)
	}

	// cho goroutine sẵn sàng
	time.Sleep(100 * time.Millisecond)
	close(start)

	wg.Wait()

	success := 0
	fail := 0

	for _, err := range results {
		if err == nil {
			success++
		} else {
			log.Println("booking failed:", err)
			fail++
		}
	}

	if success != 1 {
		t.Fatalf("Concurrency broken: expected 1 success, got %d", success)
	}

	log.Println("✅ concurrency control works correctly")
}
