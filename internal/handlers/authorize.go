package handlers

import (
	"fmt"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range h.UncheckAuthMiddlewarePaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		fmt.Println("Middleware is working") ////////////////////////////////////////////
		authToken, err := getAuthTokenFromRequest(r)
		if err != nil {
			utils.JsonError(w, http.StatusUnauthorized, "error getting auth token")
			return
		}

		fingerprint := utils.GetRequestFingerprint(r)

		uad, err := jwt.ParseToken(authToken)
		if err != nil {
			utils.JsonError(w, http.StatusUnauthorized, err.Error())
			return
		}

		fmt.Println(*uad, fingerprint)

		next.ServeHTTP(w, r)
	})
}

func getAuthTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("glts-auth-token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func getRefreshTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("glts-refresh-token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
