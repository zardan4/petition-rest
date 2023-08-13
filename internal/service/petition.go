package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zardan4/petition-audit-grpc/pkg/core/audit"
	petitions "github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
)

type PetitionService struct {
	repo *repository.Repository

	auditClient AuditClient
}

func NewPetitionService(repo *repository.Repository, auditClient AuditClient) *PetitionService {
	return &PetitionService{
		repo: repo,

		auditClient: auditClient,
	}
}

func (s *PetitionService) CreatePetition(ctx context.Context, title, text string, authorId int) (int, error) {
	petitionId, err := s.repo.Petition.CreatePetition(title, text, authorId)
	if err != nil {
		return petitionId, err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_PETITION,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log petition creation: %s", err.Error())
	}
	return petitionId, nil
}

func (s *PetitionService) GetAllPetitions() ([]petitions.Petition, error) {
	return s.repo.Petition.GetAllPetitions()
}

func (s *PetitionService) GetPetition(ctx context.Context, petitionId int) (petitions.Petition, error) {
	petition, err := s.repo.Petition.GetPetition(petitionId)
	if err != nil {
		return petition, err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_PETITION,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log petition getting: %s", err.Error())
	}
	return petition, nil
}

func (s *PetitionService) DeletePetition(ctx context.Context, petitionId, userId int) error {
	err := s.repo.Petition.DeletePetition(petitionId, userId)
	if err != nil {
		return err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_PETITION,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log petition deleting: %s", err.Error())
	}
	return nil
}

func (s *PetitionService) UpdatePetition(ctx context.Context, updatedPetition petitions.UpdatePetitionInput, petitionId, userId int) error {
	err := s.repo.Petition.UpdatePetition(updatedPetition, petitionId, userId)
	if err != nil {
		return err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_UPDATE,
		Entity:    audit.ENTITY_PETITION,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log petition updating: %s", err.Error())
	}
	return nil
}
