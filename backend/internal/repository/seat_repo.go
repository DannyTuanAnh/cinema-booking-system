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
	BookSeats(ctx context.Context, tx *sql.Tx, seats []int) error
	LockSeats(ctx context.Context, tx *sql.Tx, seats []int) error
	CountSeatsForUpdate(ctx context.Context, tx *sql.Tx, seats []int) (int, error)
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

func (s *seatRepo) BookSeats(ctx context.Context, tx *sql.Tx, seats []int) error {

	_, err := tx.ExecContext(
		ctx,
		`
            UPDATE seats
            SET status = $1
            WHERE seat_id = ANY($2);
        `,
		model.SeatStatusBooked,
		pq.Array(seats),
	)
	if err != nil {
		return fmt.Errorf("failed to update seats: %w", err)
	}

	return nil
}

func (s *seatRepo) LockSeats(ctx context.Context, tx *sql.Tx, seats []int) error {

	_, err := tx.ExecContext(
		ctx,
		`
			select seat_id
			from seats
			where seat_id = any($1)
			for update;
		`,
		pq.Array(seats),
	)

	if err != nil {
		return fmt.Errorf("failed to lock seats: %w", err)
	}

	return nil
}

func (s *seatRepo) CountSeatsForUpdate(ctx context.Context, tx *sql.Tx, seats []int) (int, error) {

	var count int

	err := tx.QueryRowContext(
		ctx,
		`
			select count(*) as cnt
			from seats
			where seat_id = any($1) and status = $2;
		`,
		pq.Array(seats),
		model.SeatStatusAvailable,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count seats: %w", err)
	}

	return count, nil
}
