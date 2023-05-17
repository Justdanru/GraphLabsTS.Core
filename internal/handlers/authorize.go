package handlers

import (
	"context"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
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

		uad, err := processAuthToken(r)
		if err == ErrNoAuthToken {
			// Check refresh token and update auth token
		}
		if err != nil {
			redirectUnathorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtxKey, uad.Id)
		ctx = context.WithValue(ctx, roleCodeCtxKey, uad.RoleCode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func processAuthToken(r *http.Request) (*models.UserAuthData, error) {
	authToken, err := getAuthTokenFromRequest(r)
	if err != nil {
		return nil, ErrNoAuthToken
	}

	fingerprint := utils.GetRequestFingerprint(r)

	uad, err := jwt.ParseToken(authToken)
	if err != nil {
		return nil, ErrParsingToken
	}

	if uad.Fingerprint != fingerprint {
		return nil, ErrDiffFingerprint
	}

	return uad, nil
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

func redirectUnathorized(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
