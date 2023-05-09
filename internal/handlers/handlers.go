package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"golang.org/x/crypto/sha3"
	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"
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
		jsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

// TODO Возможно, нужен рефакторинг
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

	var isValid bool
	isValid, err = validators.IsLoginValid(userCredentials.Login)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !isValid {
		jsonError(w, http.StatusBadRequest, "wrong login format")
		return
	}
	isValid, err = validators.IsPasswordValid(userCredentials.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !isValid {
		jsonError(w, http.StatusBadRequest, "wrong password format")
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
	uad.Fingerprint = getFingerprint(r)

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

	refreshSession := &models.RefreshSession{
		RefreshToken: refreshToken,
		Fingerprint:  uad.Fingerprint,
		UserId:       uad.Id,
	}

	sessionsCount, err := h.Repo.GetRefreshSessionsCountByUserId(uad.Id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if sessionsCount > 0 {
		err = h.Repo.DeleteAllRefreshSessionsByUserId(uad.Id)
		fmt.Printf("Found sessions: %d\n", sessionsCount)
		fmt.Println("Deleted refresh session")
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Repo.AddRefreshSession(refreshSession)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"status":        http.StatusOK,
		"auth_token":    authToken,
		"refresh_token": refreshToken,
	})

	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(map[string]interface{}{
		"error": msg,
	})
	w.Write(resp)
}

func getFingerprint(r *http.Request) string {
	headers := ""
	headers += r.Header.Get("X-Forwarded-For") + " : "
	headers += r.Header.Get("User-Agent") + ":"
	headers += r.Header.Get("Accept-Language")
	return fmt.Sprintf("%x", sha3.Sum256([]byte(headers)))
}
