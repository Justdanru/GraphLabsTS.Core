package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"
	"graphlabsts.core/internal/utils"
	"graphlabsts.core/internal/validators"

	_ "github.com/go-sql-driver/mysql"
)

type Handler struct {
	Tmpl *template.Template
	Repo repo.Repo
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "loginPage", nil)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

// TODO Возможно, нужен рефакторинг
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
		Name:     "glts_auth_token",
		Value:    authToken,
		Expires:  time.Now().Add(jwt.AUTH_TOKEN_DURATION_MINUTES * time.Minute),
		HttpOnly: true,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "glts_refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(jwt.REFRESH_TOKEN_DURATION_HOURS * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, authTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
	fmt.Println("Cookies were set")

	resp, _ := json.Marshal(map[string]interface{}{
		"url": "/profile",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "profilePage", nil)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}
