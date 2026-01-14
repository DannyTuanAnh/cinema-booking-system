package show_clients

import "time"

type GetShowResponse struct {
	ShowID     int       `json:"show_id"`
	MovieTitle string    `json:"movie_title"`
	MovieID    int       `json:"movie_id"`
	ShowTime   time.Time `json:"show_time"`
}

type ShowClient interface {
	GetShowByMovieID(movieID int) ([]*GetShowResponse, error)
}
