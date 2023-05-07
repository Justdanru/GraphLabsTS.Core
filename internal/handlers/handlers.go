package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"

	_ "github.com/go-sql-driver/mysql"
)

type Handler struct {
	Tmpl *template.Template
	Repo repo.Repo
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "loginPage", nil)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	userCredentials := &models.UserCredentials{}
	err = json.Unmarshal(body, userCredentials)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	uad, err := h.Repo.Authenticate(userCredentials)
	if err == repo.ErrNoUser {
		jsonError(w, http.StatusBadRequest, "user not found")
		return
	}
	if err == repo.ErrWrongPassword {
		jsonError(w, http.StatusBadRequest, "wrong password")
		return
	}

	authToken, err := jwt.CreateAuthToken(uad)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := jwt.CreateRefreshToken(uad)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"status":        http.StatusOK,
		"auth_token":    authToken,
		"refresh_token": refreshToken,
	})

	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "profilePage", nil)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

func jsonError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(map[string]interface{}{
		"error": msg,
	})
	w.Write(resp)
}
