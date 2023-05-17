package handlers

import (
	"html/template"

	"graphlabsts.core/internal/repo"
)

type Handler struct {
	Tmpl                       *template.Template
	Repo                       repo.Repo
	UncheckAuthMiddlewarePaths []string
}
