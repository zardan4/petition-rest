package service

import (
	"context"
	"time"

	"github.com/zardan4/petition-audit-grpc/pkg/core/audit"
	"github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
	"github.com/zardan4/petition-rest/pkg/hashing"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type Authorization interface {
	CreateUser(ctx context.Context, user core.User) (int, error)
	GenerateTokens(ctx context.Context, name, password, fingerprint string) (core.JWTPair, error)
	GenerateTokensById(userid int, fingerprint string) (core.JWTPair, error)
	ParseToken(token string) (int, error)
	CheckUserExistsById(id int) bool
	// refresh sessions
	RefreshTokens(refreshToken, fingerprint string) (core.JWTPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Petition interface {
	CreatePetition(ctx context.Context, title, text string, authorId int) (int, error)
	GetAllPetitions() ([]core.Petition, error)
	GetPetition(ctx context.Context, petitionId int) (core.Petition, error)
	DeletePetition(ctx context.Context, petitionId, userId int) error
	UpdatePetition(ctx context.Context, updatedPetition core.UpdatePetitionInput, petitionId, userId int) error
}

type Subs interface {
	GetAllSubs(petitionId int) ([]core.Sub, error)
	CreateSub(ctx context.Context, petitionId, userId int) (int, error)
	DeleteSub(ctx context.Context, petitionId, userId int) error
	CheckSignature(petitionId, userId int) (bool, error)
}

type Service struct {
	AuditClient

	Authorization
	Petition
	Subs
}

func NewService(repo *repository.Repository, hasher *hashing.SHA256Hasher, signingKey []byte, tokenTTL time.Duration, auditClient AuditClient) *Service {
	return &Service{
		AuditClient: auditClient,

		Authorization: NewAuthorizationService(repo, hasher, signingKey, tokenTTL, auditClient),
		Petition:      NewPetitionService(repo, auditClient),
		Subs:          NewSubsService(repo, auditClient),
	}
}
