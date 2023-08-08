package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/zardan4/petition-rest/internal/core"
)

func (a *AuthorizationPostgres) CountAllRefreshSessionsByUserid(userid int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(DISTINCT refresh_token) FROM %s WHERE user_id=$1", refreshSessionsTable)

	var count int

	err := a.db.QueryRow(query, userid).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (a *AuthorizationPostgres) CreateRefreshSession(userid int, fingerprint string, sessionTime time.Duration) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, fingerprint, expires_in) VALUES ($1, $2, $3) RETURNING refresh_token", refreshSessionsTable)

	var refreshToken string

	if err := a.db.QueryRow(query, userid, fingerprint, sessionTime).
		Scan(&refreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (a *AuthorizationPostgres) DeleteAllRefreshSessionsByUserId(userid int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", refreshSessionsTable)

	_, err := a.db.Exec(query, userid)

	return err
}

func (a *AuthorizationPostgres) DeleteRefreshSession(refreshToken string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token=$1", refreshSessionsTable)

	res, err := a.db.Exec(query, refreshToken)

	if affected, _ := res.RowsAffected(); affected == 0 {
		return NoRowsAffectedError
	}

	return err
}

func (a *AuthorizationPostgres) RefreshTokensAndReturnUser(refreshToken, fingerprint string) (core.User, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return core.User{}, err
	}

	// get refresh session
	var refreshSession core.RefreshSession
	query := fmt.Sprintf("SELECT * FROM %s WHERE refresh_token=$1", refreshSessionsTable)

	err = tx.Get(&refreshSession, query, refreshToken)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return core.User{}, err
	}

	// delete refresh session
	query = fmt.Sprintf("DELETE FROM %s WHERE refresh_token=$1", refreshSessionsTable)

	_, err = tx.Exec(query, refreshToken)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return core.User{}, err
	}

	// check if token is not expired
	if refreshSession.CreatedAt.Add(refreshSession.ExpiresIn).Before(time.Now()) {
		tx.Commit()                                             // nolint: errcheck
		return core.User{}, errors.New("REFRESH_TOKEN_EXPIRED") // 401
	}

	// check if fingerprint is valid
	if refreshSession.Fingerprint != fingerprint {
		tx.Commit()                                               // nolint: errcheck
		return core.User{}, errors.New("INVALID_REFRESH_SESSION") // 401
	}

	// get full user info
	query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)

	var user core.User
	err = tx.Get(&user, query, refreshSession.UserId)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return core.User{}, err
	}

	return user, tx.Commit()
}
