package movie_clients

import (
	"context"
	"time"
)

type GetMovieResponse struct {
	MovieID     int       `json:"movie_id"`
	Title       string    `json:"title"`
	Duration    int       `json:"duration"`
	Description string    `json:"description"`
	UrlImage    string    `json:"url_image"`
	Rate        float64   `json:"rate"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	Director    string    `json:"director"`
	CastList    string    `json:"cast_list"`
}

type MovieClient interface {
	GetMovieDetails(ctx context.Context) ([]*GetMovieResponse, error)
}
