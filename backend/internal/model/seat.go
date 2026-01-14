package model

const (
	SeatStatusAvailable = "available"
	SeatStatusBooked    = "booked"
)

type Seat struct {
	SeatID   int    `json:"seat_id" binding:"required,gte=1" db:"seat_id"`
	ShowID   int    `json:"show_id" db:"show_id"`
	SeatName string `json:"seat_name" db:"seat_name"`
	Status   string `json:"status" db:"status"`
}
