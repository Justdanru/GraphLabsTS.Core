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

		uad, err := checkAuthToken(r)
		if err == ErrNoAuthToken {
			// Проверить токен обновления
			// Если он отсутствует, то перенаправить на главную
			// Если токен обновления есть, то:
			//		Проверить наличие токена обновления в БД
			//  	Сравнить фингерпринт запроса с фингерпринтом в БД
			//   	Если совпадают, то обновить токен авторизации и токен обновления, записать в БД
			// 		Если фингерпринты не совпадают, то удалить все сессии пользователя из БД и перенаправить на главную
			err = h.checkRefreshToken(r)
			if err == ErrNoRefreshToken {
				redirectUnathorized(w, r)
				return
			}
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
		return nil, ErrDiffFingerprint
	}

	return uad, nil
}

func (h *Handler) checkRefreshToken(r *http.Request) error {
	refreshToken, err := getRefreshTokenFromRequest(r)
	if err != nil {
		return ErrNoRefreshToken
	}

	fingerprint := utils.GetRequestFingerprint(r)

	rs, err := h.Repo.GetRefreshSessionByToken(refreshToken)
	if err == repo.ErrNoRefreshSessions {
		return ErrNoRefreshSession
	}

	if rs.Fingerprint != fingerprint {
		return ErrDiffFingerprint
	}

	return nil
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
