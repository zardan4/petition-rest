package core

import "time"

type JWTPair struct {
	AccessToken  string
	RefreshToken string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type RefreshSession struct {
	Id           int           `db:"id"`
	UserId       int           `db:"user_id"`
	RefreshToken string        `db:"refresh_token"`
	Fingerprint  string        `db:"fingerprint"`
	ExpiresIn    time.Duration `db:"expires_in"`
	CreatedAt    time.Time     `db:"created_at"`
}
