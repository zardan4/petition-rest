package repository

import (
	"github.com/jmoiron/sqlx"
	petitions "github.com/zardan4/petition-rest/internal/core"
)

type Authorization interface {
	CreateUser(user petitions.User) (int, error)
	GetUserByName(name, password string) (petitions.User, error)
}

type Petition interface {
	CreatePetition(title, text string, authorId int) (int, error)
	GetAllPetitions() ([]petitions.Petition, error)
	GetPetition(petitionId int) (petitions.Petition, error)
	DeletePetition(petitionId, userId int) error
	UpdatePetition(petition petitions.UpdatePetitionInput, petitionId, userId int) error
}

type Subs interface {
	GetAllSubs(petitionId int) ([]petitions.Sub, error)
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
