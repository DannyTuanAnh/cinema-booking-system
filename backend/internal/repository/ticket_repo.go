package repository

import (
	"context"
	"database/sql"
	"errors"

	"cinema.com/demo/internal/model"
)

type TicketRepository interface {
	GetTicketByUserID(ctx context.Context, user_id int) ([]model.GetTicketByUserIdResponse, error)
}

type ticketRepo struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepo{db: db}
}

func (t *ticketRepo) GetTicketByUserID(ctx context.Context, user_id int) ([]model.GetTicketByUserIdResponse, error) {
	query := `SELECT 
			b.booking_id, sh.show_id, b.bookat,
			s.seat_name, m.title, sh.show_time
		FROM bookings b
		JOIN seats s ON b.seat_id = s.seat_id
		JOIN shows sh ON s.show_id = sh.show_id
		JOIN movies m ON sh.movie_id = m.movie_id
		WHERE b.user_id = $1
		ORDER BY b.bookat DESC`

	rows, err := t.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.GetTicketByUserIdResponse
	for rows.Next() {
		var ticket model.GetTicketByUserIdResponse
		err := rows.Scan(&ticket.BookingID, &ticket.ShowID, &ticket.BookAt, &ticket.SeatName, &ticket.Title, &ticket.ShowTime)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, errors.New("no tickets found for the given user ID")
	}
	return tickets, nil
}
