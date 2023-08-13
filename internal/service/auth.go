package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zardan4/petition-audit-grpc/pkg/core/audit"
	"github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
	"github.com/zardan4/petition-rest/pkg/hashing"
)

type AuthorizationService struct {
	repo   *repository.Repository
	hasher *hashing.SHA256Hasher

	HMACSecret  []byte
	tokenTTL    time.Duration
	auditClient AuditClient
}

func NewAuthorizationService(repo *repository.Repository, hasher *hashing.SHA256Hasher, signingKey []byte, tokenTTL time.Duration, auditClient AuditClient) *AuthorizationService {
	return &AuthorizationService{
		repo:        repo,
		hasher:      hasher,
		HMACSecret:  signingKey,
		tokenTTL:    tokenTTL,
		auditClient: auditClient,
	}
}

func (a *AuthorizationService) CreateUser(ctx context.Context, user core.User) (int, error) {
	hashedPassword := a.hasher.Hash(user.Password)
	user.Password = hashedPassword

	userId, err := a.repo.CreateUser(user)
	if err != nil {
		return -1, err
	}

	err = a.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_REGISTER,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log user registration: %s", err.Error())
	}

	return userId, nil
}

func (a *AuthorizationService) GenerateTokensById(userid int, fingerprint string) (core.JWTPair, error) {
	user, err := a.repo.GetUserByIdWithoutPassword(userid)
	if err != nil {
		return core.JWTPair{}, nil
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(a.tokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user_id"] = user.Id

	accessToken, err := token.SignedString(a.HMACSecret)
	if err != nil {
		return core.JWTPair{}, err
	}

	// refresh token
	countOfRefreshSessions, err := a.repo.Authorization.CountAllRefreshSessionsByUserid(user.Id)
	if err != nil {
		return core.JWTPair{}, err
	}

	// if more refresh sessions than is allowed
	if countOfRefreshSessions >= viper.GetInt("auth.maxDevidesSigned") {
		// delete all previous sessions
		err = a.repo.Authorization.DeleteAllRefreshSessionsByUserId(user.Id)
		if err != nil {
			return core.JWTPair{}, err
		}
	}

	refreshTokenTTL := viper.GetDuration("auth.refreshTTL")
	refreshToken, err := a.repo.CreateRefreshSession(user.Id, fingerprint, refreshTokenTTL)
	if err != nil {
		return core.JWTPair{}, err
	}

	return core.JWTPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,

		RefreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (a *AuthorizationService) GenerateTokens(ctx context.Context, name, password, fingerprint string) (core.JWTPair, error) {
	// access token
	hashedPassword := a.hasher.Hash(password)

	user, err := a.repo.GetUserByName(name, hashedPassword)
	if err != nil {
		return core.JWTPair{}, err
	}

	tokens, err := a.GenerateTokensById(user.Id, fingerprint)
	if err != nil {
		return core.JWTPair{}, err
	}

	a.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.Id),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log user login: %s", err.Error())
	}

	return tokens, nil
}

func (a *AuthorizationService) ParseToken(token string) (int, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return a.HMACSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !parsed.Valid || !ok {
		return 0, errors.New("invalid parsed token")
	}

	userId := claims["user_id"].(float64)

	return int(userId), nil
}

func (a *AuthorizationService) RefreshTokens(refreshToken, fingerprint string) (core.JWTPair, error) {
	user, err := a.repo.Authorization.RefreshTokensAndReturnUser(refreshToken, fingerprint)
	if err != nil {
		return core.JWTPair{}, err
	}

	return a.GenerateTokensById(user.Id, fingerprint)
}

func (a *AuthorizationService) CheckUserExistsById(id int) bool {
	_, err := a.repo.GetUserByIdWithoutPassword(id)

	return err != sql.ErrNoRows
}

func (a *AuthorizationService) Logout(ctx context.Context, refreshToken string) error {
	userId, err := a.repo.DeleteRefreshSession(refreshToken)
	if err != nil {
		return err
	}

	err = a.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGOUT,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.Fatalf("failed to log user logout: %s", err.Error())
	}
	return nil
}
