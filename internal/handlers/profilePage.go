package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/repo"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, "wrong user id")
		return
	}

	ctx, err := h.getProfilePagetContext(userId)
	if (err != nil) && (err != ErrEntityNotFound) {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}

	if err == ErrEntityNotFound {
		w.WriteHeader(http.StatusNotFound)
		err = h.Tmpl.ExecuteTemplate(w, "notFoundPage", nil)
	} else {
		err = h.Tmpl.ExecuteTemplate(w, "profilePage", ctx)
	}
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

func (h *Handler) getProfilePageContext(userId int64) (*models.ProfilePageContext, error) {
	ctx := &models.ProfilePageContext{}

	user, err := h.Repo.GetUser(userId)
	if err == repo.ErrNoSuchEntity {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	switch user.RoleCode {
	case models.AdminRoleCode:
		ctx.RoleString = "Администратор"
	case models.TeacherRoleCode:
		ctx.RoleString = "Преподаватель"
	case models.StudentRoleCode:
		ctx.RoleString = "Студент"
	}

	ctx.UserId = user.Id
	ctx.Name = user.Name
	ctx.Surname = user.Surname
	ctx.LastName = user.LastName
	ctx.TelegramId = user.TelegramId
	ctx.CreatedAt = user.CreatedAt
	ctx.UpdatedAt = user.UpdatedAt
	ctx.Groups = user.Groups
	ctx.Subjects = user.Subjects

	return ctx, nil
}
