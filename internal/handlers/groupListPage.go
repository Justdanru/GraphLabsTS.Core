package handlers

import (
	"net/http"
	"strconv"

	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) GroupListPage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		utils.JsonError(w, http.StatusBadRequest, "wrong page parameter")
		return
	}

	ctx, err := h.getGroupListPageContext(10, 10*(page-1))
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
	ctx.UserId = r.Context().Value(userIdCtxKey).(int64)
	ctx.Page = page
	ctx.PrevPage = page - 1
	ctx.NextPage = page + 1

	err = h.Tmpl.ExecuteTemplate(w, "groupListPage", ctx)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}

func (h *Handler) getGroupListPageContext(limit int64, offset int64) (*models.GroupListPageContext, error) {
	ctx := &models.GroupListPageContext{}
	var err error

	ctx.Groups, err = h.Repo.GetStudentGroups(limit, offset)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
