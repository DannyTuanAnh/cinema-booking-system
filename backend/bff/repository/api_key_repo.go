package repo

import (
	"database/sql"
	"errors"
	"time"
)

type ApiKey struct {
	ID            int64
	ClientID      string // web / mobile / facebook
	KeyHash       string
	RateLimit     int
	RateWindowSec int
	IsActive      bool
	CreatedAt     time.Time
	RevokedAt     *time.Time
}

type ApiKeyRepository interface {
	FindByHash(hash string) (*ApiKey, error)
	Insert(clientID, hash string, maxReq, winSec int) error
	ExistsByClient(clientType string) (bool, error)
}
type apiKeyRepo struct {
	db *sql.DB
}

func NewApiKeyRepo(db *sql.DB) ApiKeyRepository {
	return &apiKeyRepo{db: db}
}

func (r *apiKeyRepo) FindByHash(hash string) (*ApiKey, error) {
	var k ApiKey

	err := r.db.QueryRow(`
		SELECT id, client_id, key_hash, rate_limit, rate_window_sec, is_active, created_at, revoked_at
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
		LIMIT 1
	`, hash).Scan(
		&k.ID,
		&k.ClientID,
		&k.KeyHash,
		&k.RateLimit,
		&k.RateWindowSec,
		&k.IsActive,
		&k.CreatedAt,
		&k.RevokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("api key not found")
		}
		return nil, err
	}

	return &k, nil
}

func (r *apiKeyRepo) Insert(clientID, hash string, maxReq, winSec int) error {
	_, err := r.db.Exec(`
		INSERT INTO api_keys (client_id, key_hash, rate_limit, rate_window_sec)
		VALUES ($1, $2, $3, $4)
	`, clientID, hash, maxReq, winSec)
	return err
}

func (r *apiKeyRepo) ExistsByClient(clientType string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM api_keys WHERE client_id = $1
		)
	`, clientType).Scan(&exists)
	return exists, err
}
