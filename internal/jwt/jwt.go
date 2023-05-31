package jwt

import (
	"errors"
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

var (
	ErrWrongSigningAlg = errors.New("wrong signing alg")
	ErrNotValidToken   = errors.New("token is not valid")
)

func GetAuthTokenExpTime() time.Time {
	return time.Now().Add(AUTH_TOKEN_DURATION_MINUTES * time.Minute)
	// return time.Now().Add(5 * time.Second)
}

func GetRefreshTokenExpTime() time.Time {
	return time.Now().Add(REFRESH_TOKEN_DURATION_HOURS * time.Hour)
	// return time.Now().Add(20 * time.Second)
}

func createToken(uad *models.UserAuthData, expTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          uad.Id,
		"role":        uad.RoleCode,
		"fingerprint": uad.Fingerprint,
		"exp":         expTime,
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createAuthToken(uad *models.UserAuthData) (string, error) {
	tokenString, err := createToken(uad, GetAuthTokenExpTime())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createRefreshToken(uad *models.UserAuthData) (string, error) {
	tokenString, err := createToken(uad, GetRefreshTokenExpTime())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokenPair(uad *models.UserAuthData) (*models.TokenPair, error) {
	var err error
	tokenPair := &models.TokenPair{}

	tokenPair.AuthToken, err = createAuthToken(uad)
	if err != nil {
		return nil, err
	}

	tokenPair.RefreshToken, err = createRefreshToken(uad)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func ParseToken(tokenString string) (*models.UserAuthData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongSigningAlg
		}

		return []byte("secret JWT key"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrNotValidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrNotValidToken
	}

	uad := &models.UserAuthData{
		Id:          int64(claims["id"].(float64)),
		RoleCode:    int64(claims["role"].(float64)),
		Fingerprint: claims["fingerprint"].(string),
	}

	return uad, nil
}
