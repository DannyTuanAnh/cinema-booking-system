package domain

import "errors"

var (
	ErrSeatAlreadyBooked = errors.New("one or more seats are already booked")
)
