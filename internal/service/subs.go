package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zardan4/petition-audit-rabbitmq/pkg/core/audit"
	petitions "github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
)

type SubsService struct {
	repo *repository.Repository

	auditClient AuditClient
}

func NewSubsService(repo *repository.Repository, auditClient AuditClient) *SubsService {
	return &SubsService{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (s *SubsService) GetAllSubs(petitionId int) ([]petitions.Sub, error) {
	return s.repo.Subs.GetAllSubs(petitionId)
}

func (s *SubsService) CreateSub(ctx context.Context, petitionId, userId int) (int, error) {
	subId, err := s.repo.Subs.CreateSub(petitionId, userId)
	if err != nil {
		return subId, err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_SIGNATURE,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log signature creating: %s", err.Error())
	}
	return subId, nil
}

func (s *SubsService) DeleteSub(ctx context.Context, petitionId, userId int) error {
	err := s.repo.Subs.DeleteSub(petitionId, userId)
	if err != nil {
		return err
	}
	err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_SIGNATURE,
		EntityID:  int64(petitionId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log signature deleting: %s", err.Error())
	}
	return nil
}

func (s *SubsService) CheckSignature(petitionId, userId int) (bool, error) {
	return s.repo.CheckSignature(petitionId, userId)
}
