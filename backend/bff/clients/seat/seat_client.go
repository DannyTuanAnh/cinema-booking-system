package seat_clients

type GetSeatResponse struct {
	SeatID   int    `json:"seat_id"`
	ShowID   int    `json:"show_id"`
	SeatName string `json:"seat_name"`
	Status   string `json:"status"`
}

type SeatClient interface {
	GetSeatByShowID(showID int) ([]*GetSeatResponse, error)
}
