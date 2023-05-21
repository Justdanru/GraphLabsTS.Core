package handlers

import (
	"errors"
	"html/template"

	"graphlabsts.core/internal/repo"
)

type Handler struct {
	Tmpl                       *template.Template
	Repo                       repo.Repo
	UncheckAuthMiddlewarePaths []string
}

type ctxKey string

const (
	userIdCtxKey                  ctxKey = "UserId"
	roleCodeCtxKey                ctxKey = "RoleCode"
	MAX_REFRESH_SESSIONS_PER_USER        = 1
)

var (
	ErrWrongLoginFormat    = errors.New("wrong login format")
	ErrWrongPasswordFormat = errors.New("wrong password format")
	ErrNoAuthToken         = errors.New("no auth token")
	ErrNoRefreshToken      = errors.New("no refresh token")
	ErrParsingToken        = errors.New("error parsing token")
	ErrDiffFingerprint     = errors.New("different fingerprint")
	ErrNoRefreshSession    = errors.New("no refresh session")
)
