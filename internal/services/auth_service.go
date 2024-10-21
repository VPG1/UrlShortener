package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"url-shortener/internal/Logger"
	"url-shortener/internal/entities"
	_ "url-shortener/pkg/hasher"
)

const signingKey = "AllYourBase"

type PasswordHasher interface {
	GenerateHash(password string) string
}

type UserStorage interface {
	AddUser(name string, userName string, passwordHash string) (*entities.User, error)
	GetUser(userName string, passwordHash string) (*entities.User, error)
}

type AuthService struct {
	AliasLen       int
	Storage        UserStorage
	Logger         Logger.Logger
	PasswordHasher PasswordHasher
}

func NewAuthService(storage UserStorage, logger Logger.Logger, hasher PasswordHasher) *AuthService {
	return &AuthService{Storage: storage, Logger: logger, PasswordHasher: hasher}
}

func (as *AuthService) CreateUser(name string, userName string, password string) (*entities.User, error) {
	hashedPassword := as.PasswordHasher.GenerateHash(password)

	user, err := as.Storage.AddUser(name, userName, hashedPassword)
	if err != nil {
		as.Logger.Error("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId uint64 `json:"user_id"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.Storage.GetUser(username, s.PasswordHasher.GenerateHash(password))
	if err != nil {
		return "", err
	}

	_ = user
	// TODO: move ttl to config
	claims := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: move signingKey to envs
	return token.SignedString([]byte(signingKey))
}

func (as *AuthService) ParseToken(token string) (uint64, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			as.Logger.Error("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		as.Logger.Error("Error parsing token: %v", err)
		return 0, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		as.Logger.Error("Error parsing token claims")
		return 0, fmt.Errorf("invalid token claims")
	}

	return claims.UserId, nil
}
