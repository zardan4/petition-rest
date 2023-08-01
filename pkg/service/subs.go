package service

import (
	petitions "github.com/zardan4/petition-rest"
	"github.com/zardan4/petition-rest/pkg/repository"
)

type SubsService struct {
	repo *repository.Repository
}

func NewSubsService(repo *repository.Repository) *SubsService {
	return &SubsService{
		repo: repo,
	}
}

func (s *SubsService) GetAllSubs(petitionId int) ([]petitions.Sub, error) {
	return s.repo.Subs.GetAllSubs(petitionId)
}

func (s *SubsService) CreateSub(petitionId, userId int) (int, error) {
	return s.repo.Subs.CreateSub(petitionId, userId)
}

func (s *SubsService) DeleteSub(petitionId, userId int) error {
	return s.repo.Subs.DeleteSub(petitionId, userId)
}

func (s *SubsService) CheckSignature(petitionId, userId int) (bool, error) {
	return s.repo.CheckSignature(petitionId, userId)
}
