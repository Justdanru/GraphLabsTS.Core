package handlers

import (
	"net/http"

	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	ctx, err := h.getProfilePagetContext(r)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "profilePage", ctx)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

func (h *Handler) getProfilePagetContext(r *http.Request) (*models.ProfilePageContext, error) {
	ctx := &models.ProfilePageContext{}

	switch r.Context().Value(roleCodeCtxKey) {
	case models.AdminRoleCode:
		ctx.RoleString = "Администратор"
	case models.TeacherRoleCode:
		ctx.RoleString = "Преподаватель"
	case models.StudentRoleCode:
		ctx.RoleString = "Студент"
	}

	user, err := h.Repo.GetUser(r.Context().Value(userIdCtxKey).(int64))
	if err != nil {
		return nil, err
	}

	ctx.Name = user.Name
	ctx.Surname = user.Surname
	ctx.LastName = user.LastName
	ctx.TelegramId = user.TelegramId

	return ctx, nil
}
