package repo

import (
	"errors"

	"graphlabsts.core/internal/models"
)

type Repo interface {
	Connect(dsn string) error

	Authenticate(userCredentials *models.UserCredentials) (*models.UserAuthData, error)

	AddRefreshSession(rs *models.RefreshSession) error
	GetRefreshSessionsCountByUserId(userId int64) (int, error)
	GetRefreshSessionByToken(token string) (*models.RefreshSession, error)
	UpdateRefreshSession(rs *models.RefreshSession, oldRefreshToken string) error
	DeleteAllRefreshSessionsByUserId(userId int64) error

	GetUser(userId int64) (*models.User, error)

	GetStudentGroups(limit int64, offset int64) ([]*models.Group, error)

	GetSubjects(limit int64, offset int64) ([]*models.Subject, error)
}

var (
	ErrConnectingDB      = errors.New("error connecting db")
	ErrNoUser            = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
	ErrNoRefreshSessions = errors.New("refresh session not found")
	ErrNoSuchEntity      = errors.New("entity not found")
)
