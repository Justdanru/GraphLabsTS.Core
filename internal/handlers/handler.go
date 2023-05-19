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

var (
	ErrNoAuthToken      = errors.New("no auth token")
	ErrNoRefreshToken   = errors.New("no refresh token")
	ErrParsingToken     = errors.New("error parsing token")
	ErrDiffFingerprint  = errors.New("different fingerprint")
	ErrNoRefreshSession = errors.New("no refresh session")
)

type ctxKey string

const (
	userIdCtxKey   ctxKey = "UserId"
	roleCodeCtxKey ctxKey = "RoleCode"
)
