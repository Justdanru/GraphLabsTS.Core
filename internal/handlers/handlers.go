package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
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
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCredentials := &models.UserCredentials{}
	err = json.Unmarshal(body, userCredentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uad, err := h.Repo.Authenticate(userCredentials)
	if err == repo.ErrNoUser {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	if err == repo.ErrWrongPassword {
		http.Error(w, "wrong password", http.StatusBadRequest)
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
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func jsonError(w io.Writer, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.Write(resp)
}
