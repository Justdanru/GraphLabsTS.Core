package repo

import (
	"errors"

	"graphlabsts.core/internal/models"
)

type Repo interface {
	Connect(dsn string) error
	Authenticate(userCredentials *models.UserCredentials) (*models.UserAuthData, error)
	GetRefreshSessionsCountByUserId(userId int64) (int, error)
	GetRefreshSessionByToken(token string) (*models.RefreshSession, error)
	AddRefreshSession(rs *models.RefreshSession) error
	DeleteAllRefreshSessionsByUserId(userId int64) error
}

var (
	ErrConnectingDB      = errors.New("error connecting db")
	ErrNoUser            = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
	ErrNoRefreshSessions = errors.New("refresh session not found")
)
