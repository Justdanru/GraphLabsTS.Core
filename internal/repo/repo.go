package repo

import "graphlabsts.core/internal/models"

type Repo interface {
	Connect(dsn string) error
	Authenticate(userCredentials *models.UserCredentials) (*models.UserAuthData, error)
	GetRefreshSessionsCountByUserId(userId int64) (int, error)
	AddRefreshSession(rs *models.RefreshSession) error
	DeleteAllRefreshSessionsByUserId(userId int64) error
}
