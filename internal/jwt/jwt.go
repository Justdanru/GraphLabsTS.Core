package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"graphlabsts.core/internal/models"
)

var jwtSecretKey = []byte("secret JWT key")

// TODO Возможно, стоит перенести это в переменные окружения
const (
	AUTH_TOKEN_DURATION_MINUTES  = 5
	REFRESH_TOKEN_DURATION_HOURS = 5
)

func createToken(userID int64, userRoleCode int64, expTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": userRoleCode,
		"exp":  expTime,
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateAuthToken(uad *models.UserAuthData) (string, error) {
	tokenString, err := createToken(uad.Id, uad.RoleCode, time.Now().Add(AUTH_TOKEN_DURATION_MINUTES*time.Minute))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(uad *models.UserAuthData) (string, error) {
	tokenString, err := createToken(uad.Id, uad.RoleCode, time.Now().Add(REFRESH_TOKEN_DURATION_HOURS*time.Hour))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
