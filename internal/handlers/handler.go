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

const (
	MAX_REFRESH_SESSIONS_PER_USER = 1
)

var (
	ErrWrongLoginFormat    = errors.New("wrong login format")
	ErrWrongPasswordFormat = errors.New("wrong password format")
)
