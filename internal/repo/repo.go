package repo

import "graphlabsts.core/internal/models"

type Repo interface {
	Connect(dsn string) error
	Authenticate(login string, password string) (*models.UserAuthData, error)
}
