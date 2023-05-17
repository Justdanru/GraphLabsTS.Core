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
	ErrNoAuthToken     = errors.New("no auth token")
	ErrParsingToken    = errors.New("error parsing token")
	ErrDiffFingerprint = errors.New("different fingerprint")
)

type ctxKey string

const (
	userIdCtxKey   ctxKey = "UserId"
	roleCodeCtxKey ctxKey = "RoleCode"
)
