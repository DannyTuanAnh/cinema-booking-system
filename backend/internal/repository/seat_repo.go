package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cinema.com/demo/internal/model"
	"github.com/lib/pq"
)

type SeatRepository interface {
	GetSeatByShowID(ctx context.Context, show_id int) ([]model.Seat, error)
	BookSeats(ctx context.Context, userID int64, seats []int) error
}

type seatRepo struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepo{db: db}
}

func (s *seatRepo) GetSeatByShowID(ctx context.Context, show_id int) ([]model.Seat, error) {
	query := `SELECT seat_id, show_id, seat_name, status FROM seats WHERE show_id = $1 ORDER BY seat_name`

	rows, err := s.db.QueryContext(ctx, query, show_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []model.Seat
	for rows.Next() {
		var seat model.Seat
		err := rows.Scan(&seat.SeatID, &seat.ShowID, &seat.SeatName, &seat.Status)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(seats) == 0 {
		return nil, errors.New("no seats found for the given show ID")
	}

	return seats, nil
}

func (s *seatRepo) BookSeats(ctx context.Context, userID int64, seats []int) error {

	// Hàm có sẵn trong database để đặt chỗ
	_, err := s.db.ExecContext(
		ctx,
		"select book_seats($1, $2);",
		userID,
		pq.Array(seats),
	)
	if err != nil {
		return fmt.Errorf("failed to update seats: %w", err)
	}

	return nil
}
