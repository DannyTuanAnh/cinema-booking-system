package model

import "time"

type Show struct {
	ShowID     int       `json:"show_id" db:"show_id"`
	MovieTitle string    `json:"movie_title" db:"title"`
	MovieID    int       `json:"movie_id" db:"movie_id"`
	ShowTime   time.Time `json:"show_time" db:"show_time"`
}
