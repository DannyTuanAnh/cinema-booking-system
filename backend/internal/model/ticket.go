package model

import "time"

type GetTicketByUserIdResponse struct {
	BookingID int       `json:"booking_id" db:"booking_id"`
	ShowID    int       `json:"show_id" db:"show_id"`
	Title     string    `json:"title" db:"title"`
	SeatName  string    `json:"seat_name" db:"seat_name"`
	ShowTime  time.Time `json:"show_time" db:"show_time"`
	BookAt    time.Time `json:"book_at" db:"bookat"`
}
