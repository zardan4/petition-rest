package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable           = "users"
	petitionsTable       = "petitions_lists"
	subsTable            = "subs_items"
	usersPetitionsTable  = "users_lists"
	petitionsSubsTable   = "petitions_items"
	refreshSessionsTable = "refresh_sessions"
)

type PostgresConfig struct {
	Host string
	Port string

	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg PostgresConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
