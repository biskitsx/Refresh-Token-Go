package service

import (
	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/model"
)

type SessionService interface {
	CreateSession(userId uint) *model.Session
}

type sessionService struct {
	container container.Container
}

func NewSessionService(c container.Container) SessionService {
	return &sessionService{
		container: c,
	}
}

func (service *sessionService) CreateSession(userId uint) *model.Session {
	db := service.container.GetDatabase()
	session := &model.Session{
		UserID: userId,
	}
	db.Create(session)
	return session
}