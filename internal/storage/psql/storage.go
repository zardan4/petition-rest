package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zardan4/petition-rest/internal/core"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GetUserByName(name, password string) (core.User, error)
	GetUserByIdWithoutPassword(id int) (core.User, error)

	// refresh sessions
	CreateRefreshSession(userid int, fingerprint string, sessionTimeHours time.Duration) (string, error)
	CountAllRefreshSessionsByUserid(userid int) (int, error)
	DeleteAllRefreshSessionsByUserId(userid int) error
	DeleteRefreshSession(refreshToken string) (int, error)
	RefreshTokensAndReturnUser(refreshToken, fingerprint string) (core.User, error)
}

type Petition interface {
	CreatePetition(title, text string, authorId int) (int, error)
	GetAllPetitions() ([]core.Petition, error)
	GetPetition(petitionId int) (core.Petition, error)
	DeletePetition(petitionId, userId int) error
	UpdatePetition(petition core.UpdatePetitionInput, petitionId, userId int) error
}

type Subs interface {
	GetAllSubs(petitionId int) ([]core.Sub, error)
	CreateSub(petitionId, userId int) (int, error)
	DeleteSub(petitionId, userId int) error
	CheckSignature(petitionId, userId int) (bool, error)
}

type Repository struct {
	Authorization
	Petition
	Subs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthorizationPostgres(db),
		Petition:      NewPetitionPostgres(db),
		Subs:          NewSubsPostgres(db),
	}
}
