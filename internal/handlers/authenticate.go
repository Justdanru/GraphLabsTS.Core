package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/utils"
	"graphlabsts.core/internal/validators"
)

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

	err = validateCredentials(userCredentials)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	uad, err := h.Repo.Authenticate(userCredentials)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	uad.Fingerprint = utils.GetRequestFingerprint(r)

	tokenPair, err := jwt.CreateTokenPair(uad)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
	}

	err = h.AddRefreshSession(uad, tokenPair.RefreshToken)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, err.Error())
	}

	SetSessionCookies(w, tokenPair)

	resp, _ := json.Marshal(map[string]interface{}{
		"url": "/users/" + fmt.Sprint(uad.Id),
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func validateCredentials(userCredentials *models.UserCredentials) error {
	isValid, err := validators.IsLoginValid(userCredentials.Login)
	if err != nil {
		return err
	} else if !isValid {
		return ErrWrongLoginFormat
	}

	isValid, err = validators.IsPasswordValid(userCredentials.Password)
	if err != nil {
		return err
	} else if !isValid {
		return ErrWrongPasswordFormat
	}

	return nil
}
