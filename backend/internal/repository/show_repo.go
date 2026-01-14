package repository

import (
	"context"
	"database/sql"
	"errors"

	"cinema.com/demo/internal/model"
)

type ShowRepository interface {
	GetShowByMovieID(ctx context.Context, movie_id int) ([]model.Show, error)
}

type showRepo struct {
	db *sql.DB
}

func NewShowRepository(db *sql.DB) ShowRepository {
	return &showRepo{db: db}
}

func (s *showRepo) GetShowByMovieID(ctx context.Context, movie_id int) ([]model.Show, error) {
	query := `SELECT show_id, s.movie_id, m.title, show_time 
		FROM shows s
		JOIN movies m on s.movie_id = m.movie_id
		WHERE s.movie_id = $1 ORDER BY show_time`

	rows, err := s.db.QueryContext(ctx, query, movie_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []model.Show
	for rows.Next() {
		var show model.Show
		err := rows.Scan(&show.ShowID, &show.MovieID, &show.MovieTitle, &show.ShowTime)
		if err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(shows) == 0 {
		return nil, errors.New("no shows found for the given movie ID")
	}

	return shows, nil
}
