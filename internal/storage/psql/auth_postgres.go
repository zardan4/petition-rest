package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	petitions "github.com/zardan4/petition-rest/internal/core"
)

type AuthorizationPostgres struct {
	db *sqlx.DB
}

func NewAuthorizationPostgres(db *sqlx.DB) *AuthorizationPostgres {
	return &AuthorizationPostgres{
		db: db,
	}
}

func (a *AuthorizationPostgres) CreateUser(user petitions.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, grade, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := a.db.QueryRow(query, user.Name, user.Grade, user.Password)

	var id int
	err := row.Scan(&id)

	return id, err
}

func (a *AuthorizationPostgres) GetUserByName(name, password string) (petitions.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1 AND password_hash=$2", usersTable)

	row := a.db.QueryRow(query, name, password)

	var user petitions.User
	err := row.Scan(&user.Id, &user.Name, &user.Grade, &user.Password)

	return user, err
}
