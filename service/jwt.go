package service

import (
	"fmt"

	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/dgrijalva/jwt-go"
)

type JwtService interface {
	GenerateToken(payload interface{}, expiresIn string) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	container container.Container
	secretKey string
}

func NewJwtService(c container.Container) JwtService {
	return &jwtService{
		container: c,
		secretKey: "very-secret",
	}
}

func (service *jwtService) GenerateToken(payload interface{}, expiresIn string) (string, error) {
	claims := jwt.MapClaims{
		"user": payload,
		"exp":  expiresIn,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *jwtService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(service.secretKey), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}
