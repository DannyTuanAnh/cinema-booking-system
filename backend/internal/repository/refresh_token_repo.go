package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"cinema.com/demo/internal/model"
)

type RefreshTokenRepo interface {
	FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id int64, replacedBy string) error
	Save(ctx context.Context, rt *model.RefreshToken) error
	RevokeAllByUser(ctx context.Context, userID int64) error
}

type RefreshTokenRepoPG struct {
	db *sql.DB
}

func NewRefreshTokenRepoPG(db *sql.DB) RefreshTokenRepo {
	return &RefreshTokenRepoPG{db: db}
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *RefreshTokenRepoPG) FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {

	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`, hash)

	rt := &model.RefreshToken{}
	err := row.Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.ExpiresAt,
		&rt.Revoked,
	)

	if err != nil {
		return nil, err
	}

	return rt, nil
}

func (r *RefreshTokenRepoPG) Save(ctx context.Context, rt *model.RefreshToken) error {

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens
		(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, rt.UserID, rt.TokenHash, rt.ExpiresAt)

	return err
}

func (r *RefreshTokenRepoPG) Revoke(ctx context.Context, id int64, replacedBy string) error {

	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked = true, replaced_by_hash = $1
		WHERE id = $2
	`, replacedBy, id)

	return err
}

func (r *RefreshTokenRepoPG) RevokeAllByUser(ctx context.Context, userID int64) error {

	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE user_id = $1
	`, userID)

	return err
}
