package handlers

import (
	"net/http"

	"graphlabsts.core/internal/models"
	"graphlabsts.core/internal/utils"
)

func (h *Handler) NewGroupPage(w http.ResponseWriter, r *http.Request) {
	ctx := &models.NewGroupPageContext{}
	ctx.UserId = r.Context().Value(userIdCtxKey).(int64)

	err := h.Tmpl.ExecuteTemplate(w, "newGroupPage", ctx)
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "Template error")
		return
	}
}
