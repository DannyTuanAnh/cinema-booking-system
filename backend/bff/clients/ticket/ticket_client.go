package ticket_clients

import "time"

type GetTicketByUserIdResponse struct {
	BookingID int       `json:"booking_id"`
	ShowID    int       `json:"show_id"`
	Title     string    `json:"title"`
	SeatName  string    `json:"seat_name"`
	ShowTime  time.Time `json:"show_time"`
	BookAt    time.Time `json:"book_at"`
}

type TicketClient interface {
	GetTicketByUserID(userID int) ([]*GetTicketByUserIdResponse, error)
}
