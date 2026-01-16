package book_clients

type Seat struct {
	SeatID int `json:"seat_id"`
}

type BookRequest struct {
	UserID int64 `json:"user_id"`
	Seats  []int `json:"seats"`
}

type BookClient interface {
	BookSeats(req BookRequest) error
}
