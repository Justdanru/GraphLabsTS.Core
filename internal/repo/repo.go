package repo

import "graphlabsts.core/internal/models"

type Repo interface {
	Connect(dsn string) error
	Authenticate(userCredentials *models.UserCredentials) (*models.UserAuthData, error)
}
