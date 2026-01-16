package repository

import (
	"context"
	"database/sql"
)

type BookRepository interface {
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	CreateBooking(ctx context.Context, tx *sql.Tx, userID int, seats []int) error
	SetTimeoutTx(ctx context.Context, tx *sql.Tx, timeout string) error
}

type bookRepo struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepo{db: db}
}

func (b *bookRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return b.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (b *bookRepo) SetTimeoutTx(ctx context.Context, tx *sql.Tx, timeout string) error {
	_, err := tx.ExecContext(ctx, `SET LOCAL lock_timeout = $1`, timeout)
	return err
}

func (b *bookRepo) CreateBooking(ctx context.Context, tx *sql.Tx, userID int, seats []int) error {
	for _, seatID := range seats {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO bookings (user_id, seat_id) VALUES ($1, $2)`,
			userID,
			seatID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
