package repository

import (
	"context"
	"database/sql"
	"fmt"

	"cinema.com/demo/internal/model"
)

type MovieRepository interface {
	GetMovies(ctx context.Context) ([]model.Movie, error)
}

type movieRepo struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepo{db: db}
}

func (r *movieRepo) GetMovies(ctx context.Context) ([]model.Movie, error) {
	query := `SELECT movie_id, title, duration, description, url_image, rate, genre, release_date, director, cast_list FROM movies ORDER BY title`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %v", err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var m model.Movie
		err := rows.Scan(&m.MovieID, &m.Title, &m.Duration, &m.Description, &m.UrlImage, &m.Rate, &m.Genre, &m.ReleaseDate, &m.Director, &m.CastList)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %v", err)
		}

		movies = append(movies, m)
	}

	return movies, nil
}
