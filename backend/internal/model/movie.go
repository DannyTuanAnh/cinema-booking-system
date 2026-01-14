package model

import "time"

type Movie struct {
	MovieID     int       `json:"movie_id" db:"movie_id"`
	Title       string    `json:"title" db:"title"`
	Duration    int       `json:"duration" db:"duration"`
	Description string    `json:"description" db:"description"`
	UrlImage    string    `json:"url_image" db:"url_image"`
	Rate        float32   `json:"rate" db:"rate"`
	Genre       string    `json:"genre" db:"genre"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Director    string    `json:"director" db:"director"`
	CastList    string    `json:"cast_list" db:"cast_list"`
}
