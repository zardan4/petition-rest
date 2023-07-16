package service

import (
	petitions "github.com/zardan4/petition-rest"
	"github.com/zardan4/petition-rest/pkg/repository"
)

type PetitionService struct {
	repo *repository.Repository
}

func NewPetitionService(repo *repository.Repository) *PetitionService {
	return &PetitionService{
		repo: repo,
	}
}

func (s *PetitionService) CreatePetition(title, text string, authorId int) (int, error) {
	return s.repo.Petition.CreatePetition(title, text, authorId)
}

func (s *PetitionService) GetAllPetitions() ([]petitions.Petition, error) {
	return s.repo.Petition.GetAllPetitions()
}

func (s *PetitionService) GetPetition(petitionId int) (petitions.Petition, error) {
	return s.repo.Petition.GetPetition(petitionId)
}

func (s *PetitionService) DeletePetition(petitionId, userId int) error {
	return s.repo.Petition.DeletePetition(petitionId, userId)
}

func (s *PetitionService) UpdatePetition(updatedPetition petitions.UpdatePetitionInput, petitionId, userId int) error {
	return s.repo.Petition.UpdatePetition(updatedPetition, petitionId, userId)
}
