package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"
	"graphlabsts.core/internal/utils"
	"graphlabsts.core/internal/validators"
)

// TODO Определённо нужен рефакторинг
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	userCredentials := &models.UserCredentials{}
	err = json.Unmarshal(body, userCredentials)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	var isValid bool
	isValid, err = validators.IsLoginValid(userCredentials.Login)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !isValid {
		utils.JsonError(w, http.StatusBadRequest, "wrong login format")
		return
	}
	isValid, err = validators.IsPasswordValid(userCredentials.Password)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !isValid {
		utils.JsonError(w, http.StatusBadRequest, "wrong password format")
		return
	}

	uad, err := h.Repo.Authenticate(userCredentials)
	if err == repo.ErrNoUser {
		utils.JsonError(w, http.StatusBadRequest, "user not found")
		return
	}
	if err == repo.ErrWrongPassword {
		utils.JsonError(w, http.StatusBadRequest, "wrong password")
		return
	}
	uad.Fingerprint = utils.GetRequestFingerprint(r)

	authToken, err := jwt.CreateAuthToken(uad)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := jwt.CreateRefreshToken(uad)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshSession := &models.RefreshSession{
		RefreshToken: refreshToken,
		Fingerprint:  uad.Fingerprint,
		UserId:       uad.Id,
	}

	sessionsCount, err := h.Repo.GetRefreshSessionsCountByUserId(uad.Id)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if sessionsCount > 0 {
		err = h.Repo.DeleteAllRefreshSessionsByUserId(uad.Id)
	}
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Repo.AddRefreshSession(refreshSession)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authTokenCookie := &http.Cookie{
		Name:     "glts-auth-token",
		Value:    authToken,
		Path:     "/",
		Expires:  jwt.GetAuthTokenExpTime(),
		HttpOnly: true,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "glts-refresh-token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  jwt.GetRefreshTokenExpTime(),
		HttpOnly: true,
	}

	http.SetCookie(w, authTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	resp, _ := json.Marshal(map[string]interface{}{
		"url": "/profile",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}