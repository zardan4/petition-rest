package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	petitions "github.com/zardan4/petition-rest/internal/core"
	repository "github.com/zardan4/petition-rest/internal/storage/psql"
)

const (
	_salt       = "!fKP5wTxWa/nwO6q2k3xiRru"
	_signingKey = "Ht+Jr9)XYIm?v5tnSP.5meHIZntVVZDN32+*"
	_tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthorizationService struct {
	repo *repository.Repository
}

func NewAuthorizationService(repo *repository.Repository) *AuthorizationService {
	return &AuthorizationService{
		repo: repo,
	}
}

func (a *AuthorizationService) CreateUser(user petitions.User) (int, error) {
	hashedPassword := hashPassword(user.Password)
	user.Password = hashedPassword

	return a.repo.CreateUser(user)
}

func (a *AuthorizationService) GenerateToken(name, password string) (string, error) {
	hashedPassword := hashPassword(password)

	user, err := a.repo.GetUserByName(name, hashedPassword)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(_tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(_signingKey))
}

func (a *AuthorizationService) ParseToken(token string) (int, error) {
	var claims tokenClaims

	parsed, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(_signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	if !parsed.Valid {
		return 0, errors.New("invalid parsed token")
	}

	return claims.UserId, nil
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(_salt)))
}
