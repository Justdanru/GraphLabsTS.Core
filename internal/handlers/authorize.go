package handlers

import (
	"context"
	"net/http"

	"graphlabsts.core/internal/jwt"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"
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

		uad := &models.UserAuthData{}
		var err error

		uad, err = checkAuthToken(r)
		if err == ErrNoAuthToken {
			var oldRefreshToken string
			uad, oldRefreshToken, err = h.checkRefreshToken(r)
			if (err == ErrNoRefreshToken) || (err == ErrParsingToken) {
				redirectUnathorized(w, r)
				return
			} else if err == ErrDiffFingerprint {
				_ = h.Repo.DeleteAllRefreshSessionsByUserId(uad.Id)
				redirectUnathorized(w, r)
				return
			}

			tokenPair, err := jwt.CreateTokenPair(uad)
			if err != nil {
				redirectUnathorized(w, r)
				return
			}

			err = h.UpdateRefreshSession(tokenPair.RefreshToken, uad, oldRefreshToken)
			if err != nil {
				_ = h.Repo.DeleteAllRefreshSessionsByUserId(uad.Id)
				redirectUnathorized(w, r)
				return
			}

			SetSessionCookies(w, tokenPair)
		} else if err != nil {
			redirectUnathorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtxKey, uad.Id)
		ctx = context.WithValue(ctx, roleCodeCtxKey, uad.RoleCode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkAuthToken(r *http.Request) (*models.UserAuthData, error) {
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
		return uad, ErrDiffFingerprint
	}

	return uad, nil
}

func (h *Handler) checkRefreshToken(r *http.Request) (*models.UserAuthData, string, error) {
	refreshToken, err := getRefreshTokenFromRequest(r)
	if err != nil {
		return nil, "", ErrNoRefreshToken
	}

	fingerprint := utils.GetRequestFingerprint(r)

	uad, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, "", ErrParsingToken
	}

	rs, err := h.Repo.GetRefreshSessionByToken(refreshToken)
	if err == repo.ErrNoRefreshSessions {
		return nil, "", ErrNoRefreshSession
	}

	if (rs.Fingerprint != fingerprint) || (uad.Fingerprint != fingerprint) {
		return nil, "", ErrDiffFingerprint
	}

	return uad, refreshToken, nil
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
