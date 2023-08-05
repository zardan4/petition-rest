package service

import (
	petitions "github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user petitions.User) (int, error)
	GenerateToken(name, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Petition interface {
	CreatePetition(title, text string, authorId int) (int, error)
	GetAllPetitions() ([]petitions.Petition, error)
	GetPetition(petitionId int) (petitions.Petition, error)
	DeletePetition(petitionId, userId int) error
	UpdatePetition(updatedPetition petitions.UpdatePetitionInput, petitionId, userId int) error
}

type Subs interface {
	GetAllSubs(petitionId int) ([]petitions.Sub, error)
	CreateSub(petitionId, userId int) (int, error)
	DeleteSub(petitionId, userId int) error
	CheckSignature(petitionId, userId int) (bool, error)
}

type Service struct {
	Authorization
	Petition
	Subs
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repo),
		Petition:      NewPetitionService(repo),
		Subs:          NewSubsService(repo),
	}
}