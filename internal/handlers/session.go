package handlers

import (
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
)

func (h *Handler) AddRefreshSession(uad *models.UserAuthData, refreshToken string) error {
	refreshSession := &models.RefreshSession{
		RefreshToken: refreshToken,
		Fingerprint:  uad.Fingerprint,
		UserId:       uad.Id,
	}

	sessionsCount, err := h.Repo.GetRefreshSessionsCountByUserId(uad.Id)
	if err != nil {
		return err
	}

	if sessionsCount > (MAX_REFRESH_SESSIONS_PER_USER - 1) {
		err = h.Repo.DeleteAllRefreshSessionsByUserId(uad.Id)
	}
	if err != nil {
		return err
	}

	err = h.Repo.AddRefreshSession(refreshSession)
	if err != nil {
		return err
	}

	return nil
}

func SetSessionCookies(w http.ResponseWriter, tokenPair *models.TokenPair) {
	authTokenCookie := &http.Cookie{
		Name:     "glts-auth-token",
		Value:    tokenPair.AuthToken,
		Path:     "/",
		Expires:  jwt.GetAuthTokenExpTime(),
		HttpOnly: true,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "glts-refresh-token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		Expires:  jwt.GetRefreshTokenExpTime(),
		HttpOnly: true,
	}

	http.SetCookie(w, authTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}
