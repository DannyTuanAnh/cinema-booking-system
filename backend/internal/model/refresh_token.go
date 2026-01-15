package model

import "time"

type RefreshToken struct {
	ID             int64
	UserID         int64
	TokenHash      string
	ExpiresAt      time.Time
	Revoked        bool
	ReplacedByHash *string
	CreatedAt      time.Time
}
