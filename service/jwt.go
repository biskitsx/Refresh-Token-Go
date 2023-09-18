package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/model"
	"github.com/dgrijalva/jwt-go"
)

type JwtService interface {
	GenerateToken(payload interface{}, expiresIn uint) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
	GenerateTokenBySession(session *model.Session, expiresIn uint) (string, error)
	ExtractAccessToken(decodedAccessToken *jwt.Token) (uint, string, error)
	ExtractRefreshToken(decodedRefreshToken *jwt.Token) (uint, error)
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

func (service *jwtService) GenerateToken(payload interface{}, expiresIn uint) (string, error) {
	claims := jwt.MapClaims{
		"user": payload,
		"exp":  time.Now().Add((time.Millisecond * time.Duration(expiresIn))).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type AccessTokenDto struct {
	SessionID uint   `json:"session_id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
}

func (service *jwtService) GenerateTokenBySession(session *model.Session, expiresIn uint) (string, error) {
	payload := AccessTokenDto{
		SessionID: session.ID,
		UserID:    session.UserID,
		Username:  session.Username,
	}

	claims := jwt.MapClaims{
		"user": payload,
		"exp":  time.Now().Add((time.Millisecond * time.Duration(expiresIn))).Unix(),
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

func (service *jwtService) ExtractAccessToken(decodedAccessToken *jwt.Token) (uint, string, error) {
	// Extract the payload from the token
	payload, ok := decodedAccessToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("Invalid token claims")
	}

	fmt.Println(payload)
	// Extract the user ID and username from the payload
	userPayload, ok := payload["user"].(map[string]interface{})
	if !ok {
		return 0, "", errors.New("Invalid user payload in token")
	}

	userIDFloat, ok := userPayload["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("Invalid user ID in token")
	}

	username, ok := userPayload["username"].(string)
	if !ok {
		return 0, "", errors.New("Invalid username in token")
	}

	userID := uint(userIDFloat)
	return userID, username, nil
}

func (service *jwtService) ExtractRefreshToken(decodedRefreshToken *jwt.Token) (uint, error) {
	// Extract the payload from the token
	payload, ok := decodedRefreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	// Extract the user ID and username from the payload
	userPayload, ok := payload["user"].(map[string]interface{})
	if !ok {
		return 0, errors.New("Invalid user payload in token")
	}

	sessionIdFloat, ok := userPayload["session_id"].(float64)
	if !ok {
		return 0, errors.New("Invalid user ID in token")
	}

	sessionId := uint(sessionIdFloat)
	return sessionId, nil
}
