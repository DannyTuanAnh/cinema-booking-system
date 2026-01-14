package repository

import (
	"context"
	"database/sql"
	"errors"

	"cinema.com/demo/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, full_name, email, password FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (full_name, email, password) VALUES ($1, $2, $3)`,
		user.FullName,
		user.Email,
		user.PasswordHash,
	)
	return err
}
