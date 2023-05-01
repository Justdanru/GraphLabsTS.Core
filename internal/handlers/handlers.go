package handlers

import (
	"html/template"
	"net/http"

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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	_, err := h.Repo.Authorize(r.FormValue("login"), r.FormValue("password"))
	if err == repo.ErrNoUser {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	if err == repo.ErrWrongPassword {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}
}

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "profilePage", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
