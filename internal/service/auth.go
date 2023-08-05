package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	petitions "github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
	"github.com/zardan4/petition-rest/pkg/hashing"
)

type AuthorizationService struct {
	repo   *repository.Repository
	hasher *hashing.SHA256Hasher

	HMACSecret []byte
	tokenTTL   time.Duration
}

func NewAuthorizationService(repo *repository.Repository, hasher *hashing.SHA256Hasher, signingKey []byte, tokenTTL time.Duration) *AuthorizationService {
	return &AuthorizationService{
		repo:       repo,
		hasher:     hasher,
		HMACSecret: signingKey,
		tokenTTL:   tokenTTL,
	}
}

func (a *AuthorizationService) CreateUser(user petitions.User) (int, error) {
	hashedPassword := a.hasher.Hash(user.Password)
	user.Password = hashedPassword

	return a.repo.CreateUser(user)
}

func (a *AuthorizationService) GenerateToken(name, password string) (string, error) {
	hashedPassword := a.hasher.Hash(password)

	user, err := a.repo.GetUserByName(name, hashedPassword)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(a.tokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user_id"] = user.Id

	return token.SignedString(a.HMACSecret)
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

func (a *AuthorizationService) CheckUserExistsById(id int) bool {
	_, err := a.repo.GetUserByIdWithoutPassword(id)

	return err != sql.ErrNoRows
}
